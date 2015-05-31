var env = process.env.NODE_ENV || "development";

module.exports = function(app) {
    //controller list
    // var bottle = require("../app/controllers/bottle")(app);
    // var store = require("../app/controllers/store")(app);
    // var wine = require("../app/controllers/wine")(app);
    var home = require("../app/controllers/home")(app);

    //routes list

    //bottle routes
    // app.get("/bottle/", bottle.index);
    // app.post("/bottle/", bottle.create);

    // app.param("bottleId", bottle.load);
    // app.get("/bottle/:bottleId", bottle.show);
    // app.put("/bottle/:bottleId", bottle.update);
    // app.delete("/bottle/:bottleId", bottle.destroy);

    // //store routes
    // app.get("/store/", store.index);
    // app.post("/store/", store.create);

    // app.param("storeId", store.load);
    // app.get("/store/:storeId", store.show);
    // app.put("/store/:storeId", store.update);
    // app.delete("/store/:storeId", store.destroy);

    // //wine routes
    // app.get("/wine/", wine.index);
    // app.post("/wine/", wineType.create);

    // app.param("wineId", wineId.load);
    // app.get("/wine/:wineId", wine.show);
    // app.put("/wine/:wineId", wineType.update);
    // app.delete("/wine/:wineId", wineType.destroy);

    //home routes
    app.get("/", home.index);
    app.get("/about", home.about);

};