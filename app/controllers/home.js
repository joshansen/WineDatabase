module.exports = function(app){
	//list model variables

	exports.index = function(req, res){
		res.render("home/index", {
	        title: "Home Page",
	        desc: "This is the home page description."
	    });
	};
	exports.about = function(req, res){
		res.render("home/about", {
	        title: "About The Site",
	        desc: "This is the about page."
	    });
	};
	return exports;
};