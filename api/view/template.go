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

        .alert-body {
            padding: 10px;
            margin: -15px;
            border-top: 1px solid #ccc;
            background: #fcf8e3;
        }

        .meta {
            margin: -15px -15px 0;
        }
        .meta li:hover {
            background: #ffc107;
        }

        .paginator {
            margin: 20px 0;
        }

    </style>
</head>

<body>
<div class="container-fluid" style="max-width: 1000px">
    <h2>{{ .Group.UpdatedAt | datetime }} - {{ .Group.Rule.Name }}</h2>
    <div class="alert alert-info" role="alert">
        <p>ID：{{ .Group.ID.Hex }}</p>
        <p>规则ID：{{ .Group.Rule.ID.Hex }}</p>
        <p>规则：{{ .Group.Rule.Rule }}</p>
        <p>数量：{{ .MessageCount }}</p>
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
            <div class="panel-heading"><span class="panel-seq">#{{ $msg.SeqNum }}</span> <b>{{ datetime $msg.CreatedAt }}</b></div>
            <div class="panel-body table-responsive">
                <div class="meta">
                    <table class="table table-striped">
                        {{ if len $msg.Tags | gt 0 }}
                            <tr>
                                <th>标签</th>
                                <td>
                                    {{ range $i, $m := $msg.Tags }}
                                        <span class="label label-info">{{ $m }}</span>
                                    {{ end }}
                                </td>
                            </tr>
                        {{ end }}
                        <tr>
                            <th>来源</th>
                            <td>{{ $msg.Origin }}</td>
                        </tr>
                        {{ range $i, $m := sort_map_human $msg.Meta }}
                            <tr class="adanos-can-fold">
                                <th>{{ $m.Key }}</th>
                                <td><pre style="margin: 0; padding: 0; line-height: 1.5;">{{ format "%v" $m.Value | remove_empty_line }}</pre></td>
                            </tr>
                        {{ end }}
                        <tr>
                            <td colspan="2" style="text-align: center">
                                <a href="#" class="more-btn">更多</a>
                            </td>
                        </tr>
                    </table>
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
        $('table:lt(2)').find('a.more-btn').parents('tr').hide();
        $('table:gt(1)').find('tr.adanos-can-fold:gt(3)').hide();
        $('a.more-btn').on('click', function (e) {
            e.preventDefault();

            $(this).parents('table').find('tr.adanos-can-fold:gt(3)').show();
            $(this).parents('tr').hide();
        })
    });
</script>
</body>
</html>`
