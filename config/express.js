
var express = require("express");
var env = process.env.NODE_ENV || "development";
var swig = require("swig");
var path = require("path");
var rootPath = path.resolve(__dirname + "../..");

module.exports = function(app) {
    var CDN = require("express-cdn")(app, {
        publicDir: rootPath + "/public",
        viewsDir: rootPath + "/app/views",
        extensions: [".swig"],
        domain: process.env.S3_STATIC_BUCKET,
        bucket: process.env.S3_STATIC_BUCKET,
        key: process.env.S3_KEY,
        secret: process.env.S3_SECRET,
        ssl: false,
        production: env === "production"
    });

    app.set("showStackError", true);

    app.use(express.static(rootPath + "/public"));

    app.engine("swig", swig.renderFile);

    // views config
    app.set("views", rootPath + "/app/views");
    app.set("view engine", "swig");


    // custom error handler
    // app.use(function (err, req, res, next) {
    //     if (err.message
    //         && (~err.message.indexOf("not found")
    //         || (~err.message.indexOf("Cast to ObjectId failed")))) {
    //         return next();
    //     }

    //     console.error(err.stack);
    //     res.status(500).render("500");
    // });

    // app.use(function (req, res, next) {
    //     res.status(404).render("404", { url: req.originalUrl });
    // });
};
