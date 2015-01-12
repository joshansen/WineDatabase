module.exports = function(db) {
    try {
        return db.model("WineType");
    } catch(e) {}
	var wineTypeSchema = new db.schema({
		type: String
	});

	return db.model("WineType", wineTypeSchema);
}