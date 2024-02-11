package report1

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"time"

	"github.com/taylormonacelli/ourarkansas/listen"
)

type Timestamp struct {
	VideoID string
	Label   string
	Time    int
}

func GenerateHTMLPage(records []listen.ClipboardEntry) error {
	tmplSrc := `
<!DOCTYPE html>

<!--
https://wickydesign.com/how-to-add-timestamps-to-embedded-youtube-videos/
-->

<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>

	<style>
	body {
		font-family: 'Courier New'
	}

	.grid-container {
        display: grid;
        grid-template-columns: 60px 90px auto;
    }

    .grid-item {
        padding: 0px;
    }

	.grid-item-seconds { text-align: right; }
	.grid-item-timestamp {
		text-align: right;
		padding-left: 1ch;
	}
	.grid-item-notes {
		text-align: left;
		padding-left: 1ch;
	}
	</style>

	<script>
  var tag = document.createElement('script');
  tag.src = "https://www.youtube.com/iframe_api";
  var firstScriptTag = document.getElementsByTagName('script')[0];
  firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
  var player;
  function onYouTubeIframeAPIReady() {
    player = new YT.Player('player', {
		videoId: '{{ .VideoID }}',
    });
  }
  function setCurrentTime(slideNum) {
	var object = [{{- range $index, $element := .Records}}{{if ne $index 0}},{{end}}{{$element.Timestamp}}{{end}}];
    player.seekTo(object[slideNum]);
  }
</script>

</head>
<body style="font-family:'Courier New'">

<div class="responsive-container"><div id="player"></div></div>
<div class="grid-container">

{{range $index, $element := .Records}}
<!-- Note {{ $index }} -->
<div class="grid-item-seconds">
	<a href="{{ $element.Content }}" target="_blank">{{ $element.Timestamp | formatSeconds2}}</a>
</div>
<div class="grid-item-timestamp">
	<a href="javascript:void(0);" onclick="setCurrentTime({{ $index }})">{{ $element.Timestamp | formatSeconds}}</a>
</div>
<div class="grid-item-notes">
	{{ $element.Notes | wrapURLs }}
</div>

{{end}}

</div>
</body>
</html>
`

	funcMap := template.FuncMap{
		"formatSeconds":  formatSeconds,
		"formatSeconds2": formatSeconds2,
		"wrapURLs":       wrapURLs,
	}

	tmpl, err := template.New("timestampTemplate").Funcs(funcMap).Parse(tmplSrc)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	file, err := os.Create("clipboard.html")
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Pass the VideoID directly to the template
	data := struct {
		Records []listen.ClipboardEntry
		VideoID string
	}{
		Records: records,
		VideoID: records[0].VideoID, // Assuming all timestamps have the same VideoID
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

func RunReport1() {
	entries, err := listen.ClipboardEntriesReverseByTimestamp()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = GenerateHTMLPage(entries)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func formatSeconds(s int) string {
	d := time.Duration(s) * time.Second
	return d.String()
}

func formatSeconds2(s int) string {
	d := time.Duration(s) * time.Second
	return fmt.Sprintf("%ds", int(d.Truncate(time.Second).Seconds()))
}

func wrapURLs(input string) template.HTML {
	urlRegex := regexp.MustCompile(`(https?://\S+)`)

	result := urlRegex.ReplaceAllStringFunc(input, func(match string) string {
		return fmt.Sprintf(`<a href="%s">%s</a>`, match, match)
	})

	return template.HTML(result)
}
