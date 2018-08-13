{{define "login"}}

<!-- Login -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Login</h2>
        <div class="golfoptions">
            <form method="post">
                <que>Username</que>
                <br>
                <input type="text" name="username" placeholder=" Username">
                <br>
                <br>
                <que>Password</que>
                <br>
                <input type="password" name="password" placeholder="Password">
                <br>
                <br>
                <input type="submit">
            </form>
            <h3>No account? -
                <a href="/signup">Signup</a>
            </h3>
        </div>
    </div>
    <div class="gamespace">
    </div>
</div>
<div class="content"></div>

{{end}} {{define "signup"}}

<!-- Signup -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Sign up</h2>
        <div class="golfoptions">
            <form method="post">
                <que>Username&nbsp;</que>
                <br>
                <input type="text" name="username" placeholder=" Username">
                <br>
                <br>
                <que>Password</que>
                <br>
                <input type="password" name="password" placeholder="Password">
                <br>
                <br>
                <input type="submit">
            </form>
            <h3>Have an account? -
                <a href="/login">Login</a>
            </h3>
        </div>
    </div>
    <div class="gamespace">
    </div>
</div>
<div class="content"></div>

{{end}} {{define "profile"}}

<!-- PROFILE -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>{{.User.Username}}'s Profile</h2>
        <div class="golfoptions">
            <h4>Games: ...</h4>
            <h4>Wins: ...</h4>
            <h4>Losses: ...</h4>
            <h4>Most Used Language: ...</h4>
            <br>
            <br>
            <form method="post">
                <input type="submit" value="Log Out" style="height:25px;width:80px">
            </form>
        </div>

    </div>

</div>
<div class="content"></div>
{{end}}