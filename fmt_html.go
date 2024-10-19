// Copyright (c) 2023 Tim van der Molen <tim@kariliq.nl>
// Copyright (c) 2024 Aloïs de Gouvello <alois@outlook.fr>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tbvdm/sigtop/errio"
	"github.com/tbvdm/sigtop/signal"
)

func htmlWriteMessages(ew *errio.Writer, msgs []signal.Message) error {
  head := `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Signul</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.css">
    <style>
        section {
            margin-bottom: 1rem;
            width: 66%;
            background-color: #3b3b3bff;
            padding: 0.5rem;
            border-radius: 5px;;
        }

        .you {
            margin-left: 33%;
            background-color: #195feeff;
        }
		span {
        	font-size: 90%;
            display: flex;
            justify-content: end;
		}
    </style>
</head>

<body>
    <header>
        <h1>Conversation</h1>
    </header>

    <main>
`
  tail := `
    </main>
</body>
</html>`

	fmt.Fprintf(ew, "%s", head)
  
	for _, msg := range msgs {
		htmlWriteMessage(ew, &msg)
	}
	fmt.Fprintf(ew, "%s", tail)
	return ew.Err()
}

func htmlWriteMessage(ew *errio.Writer, msg *signal.Message) {
	if msg.IsOutgoing() {
		fmt.Fprint(ew, "<section class=\"you\">")
	} else {
		fmt.Fprintf(ew, "<section>%s:<br>", msg.Source.DisplayName())
	}
	if msg.Type != "incoming" && msg.Type != "outgoing" {
		fmt.Fprintf(ew, " [%s message]", msg.Type)
	} else {
		var details []string
		if msg.Quote != nil {
			details = append(details, fmt.Sprintf("reply to %s on %s", msg.Quote.Recipient.DisplayName(), textShortFormatTime(msg.Quote.ID)))
		}
		if len(msg.Edits) > 0 {
			details = append(details, "edited")
		}
		if len(msg.Attachments) > 0 {
			plural := ""
			if len(msg.Attachments) > 1 {
				plural = "s"
			}
			details = append(details, fmt.Sprintf("%d attachment%s", len(msg.Attachments), plural))
		}
		if len(details) > 0 {
			fmt.Fprintf(ew, " [%s]", strings.Join(details, ", "))
		}
		if msg.Body.Text != "" {
			fmt.Fprint(ew, " "+msg.Body.Text)
		}
	}


	fmt.Fprintf(ew, "<span>%s</span></section>", htmlFormatTime(msg.TimeSent))
	fmt.Fprintln(ew)
}

func htmlFormatTime(msec int64) string {
	return time.UnixMilli(msec).Format("2006-01-02 15:04")
}
