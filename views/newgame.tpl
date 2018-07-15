{{define "newgame"}}

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

				<br><br>

				<que>Max Players &nbsp;&nbsp;</que>
				<select name="maxplayers">
                        <option value="1" selected>1</option>
                        <option value="2">2</option>
                        <option value="3">3</option>
                        <option value="4">4</option>
                        <option value="5">5</option>
                        <option value="6">6</option>
                        <option value="7">7</option>
                        <option value="8">8</option>
                    </select>
				<br><br>
				<que>Difficulty &nbsp;&nbsp;</que>
				<select name="difficulty">
						<option value="beginner" selected>Beginner</option>
                        <option value="easy">Easy</option>
                        <option value="medium" selected>Medium</option>
                        <option value="hard">Hard</option>
                    </select>
				<br><br>

				<que>Game Name &nbsp;</que>
				<input type="text" name="gamename" placeholder=" name" maxlength="20">

				<br><br><br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;
				<input type="submit" value="Tee Off" style="height:25px; width:80px">

			</form>
		</div>
	</div>
</div>
{{end}}
