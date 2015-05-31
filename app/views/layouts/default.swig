<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- <link rel="icon" href="{{ CDN("/images/favicon.png", {"raw": true}) }}" type="image/x-icon"/> -->
    <title>{{ title }} - Wine Database</title>
    {% if desc %}<meta name="description" content="{{ desc }}"/>{% endif %}

    <link href="/css/bootstrap.css" rel="stylesheet" type="text/css" />
    <link href="/css/bootstrap-theme.min.css" rel="stylesheet" type="text/css" />

    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="https://oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->
</head>
<body>
    <div class="navbar navbar-default navbar-static-top">
        <div class="container">
            <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#header-navbar">
                    <span class="sr-only">Toggle Navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand" href="/">
                    Wine App</a>
            </div>

            <div id="header-navbar" class="collapse navbar-collapse">
                <ul class="nav navbar-nav">
                    <li><a href="//twitter.com/joshansen">
                        Josh's Twitter</a></li>
                </ul>
            </div>
        </div>
    </div>
    <div class="container">
        {% block content %}{% endblock %}
    </div>

    {% block script %}{% endblock %}
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
    <script src="/js/bootstrap.min.js"></script>
</body>
</html>
