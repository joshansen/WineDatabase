var mongoose = require("mongoose");


if (!process.env.MONGODB_URL) {
    throw "ENV MONGODB_URL not specified.";
}

module.exports =
    {
        mongoose: mongoose,
        schema: mongoose.Schema,

        connect: function(callback) {
            mongoose.connect(process.env.MONGODB_URL);

            mongoose.connection.on("error", function(err) {
                console.error("Connection Error:", err)
            });

            mongoose.connection.once("open", callback);
        },

        model: function(name, schema) {
            if (!schema) {
                return mongoose.model(name);
            } else {
                return mongoose.model(name, schema);
            }
        }

    };