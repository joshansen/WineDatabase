module.exports = function(db) {
    try {
        return db.model("Store");
    } catch(e) {}

    var Bottle = require("./bottle")(db);
	var storeSchema = new db.schema({

		storeName: String,

		address: String,

		city: String,

		state: String,

		zip: String,

		website: String,

		lat: Number,

		lon: Number,

		bottles: [Bottle],
	});

	return db.model("Store", storeSchema);
}