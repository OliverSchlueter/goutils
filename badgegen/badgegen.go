package badges

import "fmt"

func Generate(label, message, color string) string {
	labelWidth := len(label)*7 + 20
	messageWidth := len(message)*7 + 20
	width := labelWidth + messageWidth

	return fmt.Sprintf(`
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="20" role="img" aria-label="%s: %s">
  <linearGradient id="s" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>

  <mask id="m">
    <rect width="%d" height="20" rx="3" fill="#fff"/>
  </mask>

  <g mask="url(#m)">
    <rect width="%d" height="20" fill="#555"/>
    <rect x="%d" width="%d" height="20" fill="%s"/>
    <rect width="%d" height="20" fill="url(#s)"/>
  </g>

  <g fill="#fff" text-anchor="middle"
     font-family="Verdana, Geneva, DejaVu Sans, sans-serif"
     font-size="11">
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">%s</text>
    <text x="%d" y="14">%s</text>

    <text x="%d" y="15" fill="#010101" fill-opacity=".3">%s</text>
    <text x="%d" y="14">%s</text>
  </g>
</svg>
`,
		width, label, message,
		width,
		labelWidth,
		labelWidth, messageWidth, color,
		width,
		labelWidth/2, label,
		labelWidth/2, label,
		labelWidth+messageWidth/2, message,
		labelWidth+messageWidth/2, message,
	)
}
