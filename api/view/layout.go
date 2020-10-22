package view

var defaultLayout = `<html>
<head>
    <title>{{ .Group.Rule.Name }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
    <link href="https://cdn.bootcss.com/twitter-bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">

    <script src="https://cdn.bootcss.com/jquery/1.12.3/jquery.min.js"></script>
    <script src="https://cdn.bootcss.com/twitter-bootstrap/3.3.7/js/bootstrap.min.js"></script>

    <style type="text/css">
        .panel-body img {
            max-width: 100%;
            border-radius: 4px;
        }
        pre {
            white-space: pre-wrap;
            background-color: initial;
            border: none;
            padding-left: 0;
        }
        .json-key {
            color: #757575;
            display: inline-block;
            margin-right: 10px;
            font-weight: bold;
            word-break: break-all;
        }
        .json-value {
            color: red;
            word-break: break-all;
        }
        .panel-seq {
            margin-right: 7px;
            float: right;
            color: #676767;
        }

        .panel-body p.line:hover {
            color: #fff;
            background: #000;
        }

        .alert-body {
            padding: 10px;
            margin: 0 -15px -15px;
            border-top: 1px solid #ccc;
            background: #fcf8e3;
        }

        .meta {
            margin: -15px -15px 0;
        }
        .meta .row:hover {
            background: #e2e2e2;
        }
        .meta .row {
            padding: 5px;
            border-bottom: 1px solid #ccc;
        }
        .meta .row .label {
            line-height: 1.5;
        }
        .value-box {
            margin: 0;
            padding: 0;
            line-height: 1.5;
            font-size: 80%;
        }
        .paginator {
            margin: 20px 0;
        }
        .bs-callout {
            padding: 20px;
            margin: 20px 0;
            border: 1px solid #eee;
            border-left-width: 5px;
            border-radius: 3px;
        }
        .bs-callout-warning {
            border-left-color: #aa6708;
        }
        .bs-callout-warning h4 {
            color: #aa6708;
        }
        .bs-callout h4 {
            margin-top: 0;
            margin-bottom: 5px;
        }
    </style>
</head>

<body>
<div class="container-fluid">
    <h2>{{ .Group.Rule.Name }} {{ if eq .Group.Type "recovery" }}<span class="label label-success">恢复</span>{{ end }}</h2>
    <div class="bs-callout bs-callout-warning">
        <p>创建时间：{{ .Group.UpdatedAt | datetime "2006-01-02 15:04:05" }}</p>
        <p>规则：<code>{{ .Group.Rule.Rule }}</code></p>
        <p>事件数量：{{ .EventsCount }}</p>
    </div>
    {{--BODY--}}
</div>

<script>
$(function() {
    $('.alert-body table').addClass("table table-hover");
});
</script>

</body>
</html>
`
