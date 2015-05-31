module.exports = function(db) {
    try {
        return db.model("Wine");
    } catch(e) {}
	var Bottle = require("./bottle")(db);

	var wineSchema = new db.schema({

		created: {type: Date, "default": Date.now},

		modified: Date,

		information: String,

		imageUrl: String,

		originalImageUrl: String,

		name: String,

		brand: String,

		types: [String],

		year: Number,

		//information from bottles
		bottles: [Bottle],


	});

	wineSchema.virtual("type")
		.get(function() {
            return this.types[0];
        })
        .set(function(type) {
            if (this.types[0]) {
                this.types[0].remove();
            }
            this.types.push(type);
        });

	// wineSchema.methods{

	// }



	return db.model("Wine", wineSchema);
};