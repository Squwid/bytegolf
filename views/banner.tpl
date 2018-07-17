{{define "banner"}}

<!-- Banner in HTML -->
<div class="banner">
    <div class="logospace"> logo space</div>

    <div class="bannerlistspace"></div>

    <a href="/">
        <div class="bannerlist" a href="golf.html">
            Tee Off
        </div>
    </a>

    <div class="bannerlistspace"></div>

    <div class="bannerlist">
        Rules
    </div>
    <div class="bannerlistspace"></div>

    <a href="/leaderboards">
        <div class="bannerlist">
            Leaderboard
        </div>
    </a>
    <div class="bannerlistspace"></div>
    <a href="/login">
        <div class="bannerlist">
            {{if .LoggedIn}} {{.Name}} {{else}} Log in/Sign up {{end}}
        </div>
    </a>
    <div class="bannerlistspace"></div>
</div>

{{end}}
