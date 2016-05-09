package models

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

// Need to implement find by name, geographic distance, etc.
// We could also do this client side

type Wine struct {
	Id           bson.ObjectId `bson:"_id"`
	CreatedDate  time.Time
	ModifiedDate time.Time
	Name         string
	Winery       string
	Information  string
	VarietyID    bson.ObjectId
	Style        string
	Region       string
	//Overall stats stored below
	MaxPrice        float64
	MaxSlug         string
	MinRegularPrice float64
	MinRegularSlug  string
	MinSalePrice    float64
	MinSaleSlug     string
	AvgPrice        float64
	AvgRating       float64
	BestYears       BestYears
	LastImage       string
}

//Define the bestYears type that will be used in single.
type BestYears []int

//Define the String method on bestYears that will be used to print a list of best years with commas.
func (ys BestYears) String() string {
	stringOfYears := ""

	for i, y := range ys {
		if i == 0 {
			stringOfYears = fmt.Sprint(y)
			continue
		}
		stringOfYears = stringOfYears + ", " + fmt.Sprint(y)
	}

	return stringOfYears
}

func appendIfMissing(slice []bson.ObjectId, i bson.ObjectId) []bson.ObjectId {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (w *Wine) Save(db *mgo.Database) error {
	_, err := w.coll(db).UpsertId(w.Id, w)
	return err
}

func (w *Wine) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return w.coll(db).FindId(id).One(w)
}

func (*Wine) coll(db *mgo.Database) *mgo.Collection {
	return db.C("wine")
}

type Wines []Wine

func (w *Wines) FindAll(db *mgo.Database) error {
	return w.coll(db).Find(nil).All(w)
}

func (*Wines) coll(db *mgo.Database) *mgo.Collection {
	return db.C("wine")
}

func (w *Wine) calculateStats(db *mgo.Database) error {
	//Calculate statistics for wine on all purchases
	purchases := new(Purchases)
	if err := purchases.FindByWineID(w.Id, db); err != nil {
		return errors.New("could not find associated purchase records.")
	}

	//Initialize the variables that will be used to calculate wine statistics.
	lenPurchases := len(*purchases)
	var sumPrice float64
	var sumRating int
	var maxRating int
	var lastBought time.Time

	//Loop over all purchase records to calculate statistics.
	for _, purchase := range *purchases {
		//Update maxPrice if current maxPrice is less than purchase price.
		if w.MaxPrice < purchase.Price {
			w.MaxPrice = purchase.Price
			w.MaxSlug = purchase.Id.Hex()
		}

		//Set minRegularPrice to first nonsale purchase.
		if w.MinRegularPrice == 0 && !purchase.OnSale {
			w.MinRegularPrice = purchase.Price
			w.MinRegularSlug = purchase.Id.Hex()
		}
		//Update minRegularPrice if current minRegularPrice is greater than purchase price and the purchase wasn't on sale.
		if w.MinRegularPrice > purchase.Price && !purchase.OnSale {
			w.MinRegularPrice = purchase.Price
			w.MinRegularSlug = purchase.Id.Hex()
		}

		//Set minSalePrice to first sale purchase.
		if w.MinSalePrice == 0 && purchase.OnSale {
			w.MinSalePrice = purchase.Price
			w.MinSaleSlug = purchase.Id.Hex()
		}
		//Update minSalePrice if current minSalerice is greater than purchase price and the purchase was on sale.
		if w.MinSalePrice > purchase.Price && purchase.OnSale {
			w.MinSalePrice = purchase.Price
			w.MinSaleSlug = purchase.Id.Hex()
		}

		//Add current purchase price to sumPrice that will be used to calculate avgPrice.
		sumPrice = sumPrice + purchase.Price

		//Update maxRating if purchase rating is greater than current maxRating.
		if maxRating < purchase.Rating {
			maxRating = purchase.Rating
		}

		//Add current purchase rating to sumRating that will be used to calculate avgRating.
		sumRating = sumRating + purchase.Rating

		//Set LastImage if it is blank and ImageOriginalURL is not.
		if w.LastImage == "" && purchase.ImageOriginalURL != "" {
			w.LastImage = purchase.ImageOriginalURL
		}
		//Update lastBought if lastBought is after date purchased.
		if lastBought.Before(purchase.DatePurchased) {
			lastBought = purchase.DatePurchased
			//Update LastImage if ImageOriginalURL isn't blank.
			if purchase.ImageOriginalURL != "" {
				w.LastImage = purchase.ImageOriginalURL
			}
		}
	}

	//Create an array of unique best years
	var maxYearMatch bool
	for _, purchase := range *purchases {
		if maxRating == purchase.Rating {
			maxYearMatch = false
			for _, year := range w.BestYears {
				if year == purchase.Year {
					maxYearMatch = true
					break
				}
			}
			if !maxYearMatch {
				w.BestYears = append(w.BestYears, purchase.Year)
			}
		}
	}

	//Reverse sort the list of best years
	sort.Sort(sort.Reverse(sort.IntSlice(w.BestYears)))

	//Calculate avgPrice and avgRating.
	w.AvgPrice = sumPrice / float64(lenPurchases)
	w.AvgRating = float64(sumRating) / float64(lenPurchases)

	return w.Save(db)
}
