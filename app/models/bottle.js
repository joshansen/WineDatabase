module.exports = function(db) {
    try {
        return db.model("Bottle");
    } catch(e) {}

    var ObjectId = db.schema.Types.ObjectId;

    var bottleSchema = new db.schema({

    	created: {type: Date, "default": Date.now},

        modified: Date,

        wine: {type: ObjectId, ref: "Wine"},

    	store: {type: ObjectId, ref: "Store"},

    	buyAgain: {type: Boolean, "default": false},

    	doWeLike: {type: Boolean, "default": false},

        notes: String,

        price: Number,

        datePurchased: Date,

        dateDrank: Date,

        memoryCue: String
    });

    return db.model("Bottle", bottleSchema);
}

