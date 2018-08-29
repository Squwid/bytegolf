{{define "banner"}}

<!-- Banner in HTML -->
<div class="banner">
    <div class="logospace">Byte Golf</div>

    <div class="bannerlistspace"></div>

    <a href="/">
        <div class="bannerlist" a href="golf.html">
            Tee Off
        </div>
    </a>

    <div class="bannerlistspace"></div>

    <a href="/rules">
        <div class="bannerlist">
            Rules
        </div>
    </a>

    <div class="bannerlistspace"></div>

    <div class="bannerlist">
        Leaderboard
    </div>
    <div class="bannerlistspace"></div>
    <a href="/login">
        <div class="bannerlist">
            {{if .LoggedIn}} {{.Name}} {{else}} Log in/Sign up {{end}}
        </div>
    </a>
    <div class="bannerlistspace"></div>
</div>

{{end}} {{define "footer"}}

<!-- Footer in HTML -->
<div class="footer">
    <div class="footerspace"></div>
    <div class="footercolumn">
        <!-- First Footer Column Area -->
    </div>
    <div class="footerspace"></div>
    <div class="footercolumn">
        <!-- Second Footer Column Area -->
    </div>
    <div class="footerspace"></div>
    <div class="footercolumn">
        <!-- Second Footer Column Area -->
    </div>
    <div class="footerspace"></div>
</div>

{{end}}