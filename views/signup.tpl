 {{define "signup"}}

<!-- Signup -->
<div class="content">
    <div class="contentspace">
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
                <h3><a href="/login">Login</a></h3>
            </div>
        </div>
        <div class="gamespace">
        </div>
    </div>
</div>

{{end}}
