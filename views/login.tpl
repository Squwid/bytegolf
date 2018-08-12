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
            <h3>No account? - <a href="/signup">Signup</a></h3>
        </div>
    </div>
    <div class="gamespace">
    </div>
</div>

{{end}}

 {{define "signup"}}

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
            <h3>Have an account? - <a href="/login">Login</a></h3>
        </div>
    </div>
    <div class="gamespace">
    </div>
</div>

{{end}}
