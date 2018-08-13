{{define "showcode"}}
<!-- Show the code whether it is correct or not -->
    {{if .OverallCode.Show}}
        {{if .OverallCode.Correct}}
            <div class="correct">{{.OverallCode.Output}}<br>
            You have gotten this hole correct. Your best score is {{.OverallCode.Bytes}}
            </div>
        {{end}}
    {{end}}
    {{if .CurrentCode.Show}}
        {{if .CurrentCode.Correct}}
            <div class="correct">{{.CurrentCode.Output}}<br>
                This output is correct. Score of {{.CurrentCode.Bytes}}
            </div>
        {{else}}
            <div class="incorrect">{{.CurrentCode.Output}}<br>
                This output is incorrect. 
            </div>
        {{end}}
    {{end}}

{{end}}


{{define "leaderboards"}}

<!-- Leaderboards -->
<div class="content">
    <div class="contentspace"></div>
    <div class="lbcontainer">
        <h2>Leaderboards ({{.Game.CurrentPlayers}}/{{.Game.MaxPlayers}} players)</h2>
        <a href="/master">More Options</a>
        <div class="playercontainer">
            <div class="gamespace"></div>
            <div class="player1">
                &nbsp;&nbsp;&nbsp;&nbsp; {{.Game.Leaderboard.Winning.User.Username}} - with {{.Game.Leaderboard.Winning.HolesCorrect}} holes <br>
                &nbsp;&nbsp;&nbsp;&nbsp; {{.Game.Leaderboard.Winning.TotalScore}} bytes
            </div>
            {{range .Game.Leaderboard.OtherPlayers}}
            <div class="gamespace"></div>
            <div class="player">
                &nbsp;&nbsp;&nbsp;&nbsp; {{.User.Username}} - with {{.HolesCorrect}} holes <br>
                &nbsp;&nbsp;&nbsp;&nbsp; {{.TotalScore}} bytes
            </div>
            {{end}}
            <div class="gamespace"></div>
        </div>
    </div>
</div>

{{end}}

{{define "question"}}

<!-- Question -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <br>
        <div class="qoptions">
            <h1> {{.Hole}}. </h1>
            <br>
            <h2> {{.Question.Name}} </h2>
            <br>
            <h4>{{ .Question.Question}}</h4>
        </div>
    </div>
    <div class="gamespace">

    </div>
</div>

{{end}}


{{define "submit"}}

<!-- Submission -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Submission</h2>
        <div class="golfoptions">
            <form method="POST" enctype="multipart/form-data">
                <que>Current Game &nbsp;</que>
                <br> <h2>{{.Game.Name}}</h2>
                <br>
                <br>
                <que>Language</que>
                <select name="language">
                        <option value="java">Java</option>
                        <option value="c">C</option>
                        <option value="cpp">C++</option>
                        <option value="cpp14">C++14</option>
                        <option value="php">PHP</option>
                        <option value="python2">Python 2</option>
                        <option value="python3">Python 3</option>
                        <option value="ruby">Ruby</option>
                        <option value="go">Go</option>
                        <option value="bash">Bash</option>
                        <option value="swift">Swift</option>
                        <option value="r">R</option>
                        <option value="nodejs">NodeJS</option>
                        <option value="fsharp">F#</option>
                    </select>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <div class="gamespace"></div>
                <br>
                <br>
                <que>Select Code File</que>
                <br>
                <br>
                <input type="file" name="codefile">
                <br><br>&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;
                <input type="submit" value="Sumbit" style="height:25px; width:80px">
            </form>
            
        </div>
    </div>

</div>

{{end}}
