package main

const (
	HTML_HEAD = `
	<html>
		<head>
			<title>Swipe</title>
			
			<style>
				body {
					background-color: #fff;
					margin: 0;
					padding: 1rem 0;
				}
				img {
					margin: 1rem 2rem;
					border: 1px solid #000;
				}
			</style>
		</head>
		<body>
	`

	HTML_TAIL = `
		</body>
	</html>
	`

	HTML_IMG_HEAD = `<img src="data:image/png;base64,`
	HTML_IMG_TAIL = `" />`
)
