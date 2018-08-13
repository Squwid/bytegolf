{{define "currentgames"}}

<!-- Current Games -->
<div class="content">
    <div class="contentspace"></div>
    <div class="cgcontainer">
        <h2>Current Games</h2>
        <div class="gamecontainer">
            <div class="gamespace"></div>
            {{if .Game.Started}}
            <a href="/currentgame/1">
                <div class="game">
                    <br> &nbsp;&nbsp; {{.Game.Name}}
                    <br> &nbsp;&nbsp; Players: {{.Game.CurrentPlayers}}
                </div>
            </a>
            {{else}}
            <h3>No current games</h3>
            {{end}}
            <div class="gamespace"></div>
            <div class="gamespace"></div>
        </div>
    </div>
</div>

{{end}} {{define "newgame"}}

<!-- New Game in HTML -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Create New Game</h2>
        <div class="golfoptions">
            <form method="post">
                <que>Hole Amount &nbsp;</que>
                <input type="radio" name="holes" value="1"> 1 &nbsp; &nbsp;
                <input type="radio" name="holes" value="3" checked> 3 &nbsp; &nbsp;
                <input type="radio" name="holes" value="9"> 9 &nbsp; &nbsp;

                <br>
                <br>

                <que>Max Players &nbsp;&nbsp;</que>
                <select name="maxplayers">
                    <option value="1">1</option>
                    <option value="2">2</option>
                    <option value="3">3</option>
                    <option value="4">4</option>
                    <option value="5" selected>5</option>
                    <option value="6">6</option>
                    <option value="7">7</option>
                    <option value="8">8</option>
                </select>

                <br>
                <br>

                <que>Game Name &nbsp;</que>
                <input type="text" name="gamename" placeholder=" Name" maxlength="15">

                <br>
                <br>

                <que>Password &nbsp; &nbsp; &nbsp;</que>
                <input type="password" name="password" maxlength="15" placeholder=" Password">

                <br>
                <br>
                <br>
                <br>
                <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                &nbsp;&nbsp;
                <input type="submit" value="Tee Off" style="height:25px; width:80px">

            </form>
        </div>
    </div>
</div>


{{end}} {{define "gamemaster"}}
<!-- Game Master Controls -->
<!-- Content Information -->
{{if .GameOver}}
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>{{.Game.Name}} Results</h2>
        <div class="golfoptions">
            <h2>Winner: {{.Game.Leaderboard.Winning.User.Username}}</h2>
            <h2>Score: {{.Game.Leaderboard.Winning.TotalScore}} bytes</h2>
        </div>
    </div>
</div>
{{end}}

<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>{{.Game.Name}} Options</h2>
        <div class="golfoptions">
            <!-- <h4>ID: {{.Game.ID}}</h4> -->
            <h4>Name: {{.Game.Name}}</h4>
            <h4>Password: {{.Game.Password}}</h4>
            <h4>Current Players: {{.Game.CurrentPlayers}}</h4>
            <h4>Max Players: {{.Game.MaxPlayers}}</h4>
            <h4>Holes: {{.Game.Holes}}</h4>
            <h4>Difficulty: {{.Game.Difficulty}}</h4>
            <h4>Start Time: {{.Game.StartedTime}}</h4>
            {{if .CanEnd}}
            <form method="post">
                <br>
                <br>
                <br>
                <br>
                <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                &nbsp;&nbsp;
                <input type="submit" value="End Game" style="height:25px; width:80px">
            </form>
            {{end}}
        </div>
    </div>
</div>
<!-- Player Information -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>{{.Game.Name}} Players</h2>
        <div class="golfoptions">
            {{range .Game.Players}}
            <div class="pinfo">
                <strong>{{.User.Username}}</strong>
                <br> {{.HolesCorrect}} Holes - {{.TotalScore}} Bytes
            </div>
            {{end}}
        </div>
    </div>
</div>


{{end}}