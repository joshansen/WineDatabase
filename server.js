var express = require("express");
var env = process.env.NODE_ENV || "development";
var fs = require("fs");

// Load in configuration options
require("dotenv").load();

var db = require("./lib/db");


db.connect(function() {
    fs.readdirSync(__dirname + "/app/models").forEach(function (file) {
        if (~file.indexOf(".js")) {
            require(__dirname + "/app/models/" + file)(db);
        }
    });

    var app = express();

    // Start the app by listening on <port>
    var port = process.env.PORT;

    console.log("PORT: " + port);

    app.listen(port, function() {
        if (process.send) {
            process.send("online");
        }
    });

    process.on("message", function(message) {
        if (message === "shutdown") {
            process.exit(0);
        }
    });
});
