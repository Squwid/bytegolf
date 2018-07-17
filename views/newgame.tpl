{{define "newgame"}}

<!-- New Game in HTML -->
<div class="content">
<<<<<<< HEAD
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Create New Game</h2>
        <div class="golfoptions">
            <form method="post">
                <que>Hole Amount &nbsp;</que>
                <input type="radio" name="holes" value="one"> 1 &nbsp; &nbsp;
                <input type="radio" name="holes" value="three" checked> 3 &nbsp; &nbsp;
                <input type="radio" name="holes" value="nine"> 9 &nbsp; &nbsp;

                <br><br>

                <que>Max Players &nbsp;&nbsp;</que>
                <select name="maxplayers">
=======
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
>>>>>>> 5a8bfce48c0efcfb0e854935b827fd3e887526c4
                        <option value="1" selected>1</option>
                        <option value="2">2</option>
                        <option value="3">3</option>
                        <option value="4">4</option>
                        <option value="5">5</option>
                        <option value="6">6</option>
                        <option value="7">7</option>
                        <option value="8">8</option>
                    </select>
<<<<<<< HEAD

                <br><br>
=======
				<br><br>
				<que>Difficulty &nbsp;&nbsp;</que>
				<select name="difficulty">
						<option value="beginner" selected>Beginner</option>
                        <option value="easy">Easy</option>
                        <option value="medium" selected>Medium</option>
                        <option value="hard">Hard</option>
                    </select>
				<br><br>
>>>>>>> 5a8bfce48c0efcfb0e854935b827fd3e887526c4

                <que>Game Name &nbsp;</que>
                <input type="text" name="gamename" placeholder=" name" maxlength="20">

<<<<<<< HEAD
                <br><br>

                <que>Password &nbsp; &nbsp; &nbsp;</que>
                <input type="password" name="password" maxlength="20" placeholder=" password">

                <br><br><br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;
                <input type="submit" value="Tee Off" style="height:25px; width:80px">
=======
				<br><br><br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;
				<input type="submit" value="Tee Off" style="height:25px; width:80px">
>>>>>>> 5a8bfce48c0efcfb0e854935b827fd3e887526c4

            </form>
        </div>
    </div>
</div>

{{end}}
