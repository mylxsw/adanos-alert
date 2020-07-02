package view


var defaultMessageViewTemplate = `<html>
<head>
    <title>{{ .Group.Rule.Name }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://cdn.bootcss.com/twitter-bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">

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

        .message-key {
            min-width: 25%; 
            display: inline-block;
        }
    </style>
</head>

<body>
<div class="container-fluid" style="max-width: 1000px">
    <h2>{{ .Group.Rule.Name }}</h2>
    <div class="alert alert-warning" role="alert">
        <p>ID：{{ .Group.ID.Hex }}</p>
        <p>规则ID：{{ .Group.Rule.ID.Hex }}</p>
        <p>规则：{{ .Group.Rule.Rule }}</p>
        <p>数量：{{ .MessageCount }}</p>
    </div>
    {{ range $i, $msg := .Messages }}
        <div class="panel panel-default">
            <div class="panel-heading"><span class="panel-seq">#{{ $i }}</span> <b>{{ datetime $msg.CreatedAt }}</b></div>
            <div class="panel-body table-responsive">
                <ul style="list-style: none;padding-left: 0;">
                    <li><span class="message-key"><b style="border-bottom: 1px dashed #000000;">标签：</b></span> {{ range $i, $m := $msg.Tags }}
                            <span class="label label-info">{{ $m }}</span>
                        {{ end }}</li>
                    <li><span class="message-key"><b style="border-bottom: 1px dashed #000000;">来源：</b></span> {{ $msg.Origin }}</li>
                    {{ range $i, $m := $msg.Meta }}
                        <li><span class="message-key"><b style="border-bottom: 1px dashed #000000;">{{ $i }} :</b></span> {{ $m }}</li>
                    {{ end }}
                </ul>
                {{ range $i, $m := json_flatten $msg.Content 3 }}
                    <div class="json-line">
                        <span class="json-key">{{ $m.Key }}</span>
                        <pre class="json-value">{{ if eq $m.Key "stack" }}{{ range $j, $k := json_flatten $m.Value 1 }}
                                {{ $k.Key }} : {{ $k.Value }}{{ end }}{{ else }}{{ $m.Value }}{{ end }}</pre>
                    </div>
                {{ else }}
                    <pre>{{ json $msg.Content }}</pre>
                {{ end }}
            </div>
        </div>
    {{ end }}
</div>
<script src="https://cdn.bootcss.com/jquery/1.12.3/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/twitter-bootstrap/3.3.7/js/bootstrap.min.js"></script>
<script>
    $(function() {

    });
</script>
</body>
</html>`
