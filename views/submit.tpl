{{define "leaderboards"}}

<!-- Leaderboards -->
<div class="content">
    <div class="contentspace"></div>
    <div class="lbcontainer">
        <h2>Leaderboards</h2>
        <div class="playercontainer">
            <div class="gamespace"></div>
            <div class="player1">
                player1
            </div>
            <div class="gamespace"></div>
            <div class="player">
                player2
            </div>
            <div class="gamespace"></div>
            <div class="player">
                player3
            </div>
            <div class="gamespace"></div>
            <div class="player">
                player4
            </div>
            <div class="gamespace"></div>
            <div class="player">
                player5
            </div>
            <div class="gamespace"></div>
            <div class="player">
                player6
            </div>
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
                <br> <h2>{{.GameName}}</h2>
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
                <br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;
                <input type="submit" value="Sumbit" style="height:25px; width:80px">

            </form>
        </div>
    </div>

</div>

{{end}}
