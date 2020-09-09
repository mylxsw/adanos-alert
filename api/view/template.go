package view

var defaultMessageViewTemplate = `<html>
<head>
    <title>{{ .Group.Rule.Name }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
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
    <h2>{{ .Group.Rule.Name }}</h2>
    <div class="bs-callout bs-callout-warning">
        <p>创建时间：{{ .Group.UpdatedAt | datetime "2006-01-02 15:04:05" }}</p>
        <p>规则：<code>{{ .Group.Rule.Rule }}</code></p>
        <p>消息数量：{{ .MessageCount }}</p>
    </div>
    <div class="paginator">
        {{ if .HasPrev }}
            <a class="btn btn-info" href="{{ .Path }}?offset={{ .PrevOffset }}" role="button">上一页</a>
        {{ end }}

        {{ if .HasNext }}
            <a class="btn btn-primary" href="{{ .Path }}?offset={{ .Next }}" role="button">下一页</a>
        {{ end }}
    </div>
    {{ range $i, $msg := .Messages }}
        <div class="panel panel-default">
            <div class="panel-heading"><span class="panel-seq" id="{{ $msg.SeqNum }}"><a href="#{{ $msg.SeqNum }}">#{{ $msg.SeqNum }}</a></span> <b>{{ datetime "2006-01-02 15:04:05" $msg.CreatedAt }}</b></div>
            <div class="panel-body">
                <div class="meta">
                    <div class="container-fluid">
                        {{ if len $msg.Tags | gt 0 }}
                            <div class="row">
                                <div class="col-sm-3"><span class="label label-success">标签</span></div>
                                <div class="col-sm-9">
                                    {{ range $i, $m := $msg.Tags }}
                                        <span class="label label-info">{{ $m }}</span>
                                    {{ end }}
                                </div>
                            </div>
                        {{ end }}
                        <div class="row">
                            <div class="col-sm-3"><span class="label label-primary">来源</span></div>
                            <div class="col-sm-9">
                                {{ $msg.Origin }}
                            </div>
                        </div>
                        {{ range $i, $m := sort_map_human $msg.Meta }}
                            {{ if format "%v" $m.Value | ne "" }}
                            <div class="row adanos-can-fold">
                                <div class="col-sm-3"><span class="label label-info">{{ $m.Key }}</span></div>
                                <div class="col-sm-9">
                                    <pre class="value-box">{{ format "%v" $m.Value | remove_empty_line }}</pre>
                                </div>
                            </div>
                            {{ end }}
                        {{ end }}
                    </div>
                </div>
                <div class="alert-body">
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
        </div>
    {{ end }}

    <div class="paginator">
        {{ if .HasPrev }}
            <a class="btn btn-info" href="{{ .Path }}?offset={{ .PrevOffset }}" role="button">上一页</a>
        {{ end }}

        {{ if .HasNext }}
            <a class="btn btn-primary" href="{{ .Path }}?offset={{ .Next }}" role="button">下一页</a>
        {{ end }}
    </div>
</div>
<script src="https://cdn.bootcss.com/jquery/1.12.3/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/twitter-bootstrap/3.3.7/js/bootstrap.min.js"></script>
<script>
    $(function() {

    });
</script>
</body>
</html>`
