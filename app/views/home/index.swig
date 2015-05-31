{% extends "../layouts/default.swig" %}

{% block title %}{{Home page!}}{% endblock %}

{% block content %}

<div class="row" id="like-row">
    <div class="col-md-12">
        <p>
            Like Row
        </p>
    </div>
</div>
<div class="row" id="mixed-row">
    <div class="col-md-12">
        <p>
            Mixed Like Row
        </p>
    </div>
</div>
<div class="row" id="dont-like-row">
    <div class="col-md-12">
        <p>
            Don't Like Row
        </p>
    </div>
</div>

{% endblock %}
