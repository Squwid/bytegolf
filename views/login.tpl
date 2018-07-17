{{define "login"}}

<!-- Login -->
<div class="content">
    <div class="contentspace">
        <div class="golfcontainer">
            <h2>Login</h2>
            <div class="golfoptions">
                <form method="post">
                    <que>Username&nbsp;</que>

                    <input type="text" name="username" placeholder=" Username">
                    <br>
                    <br>

                    <que>Password&nbsp;</que>
                    <input type="password" name="password" placeholder="Password">
                    <br>
                    <br>
                    <input type="submit">
                </form>
                <h3><a href="/signup">signup</a></h3>
            </div>
        </div>
        <div class="gamespace">

        </div>
    </div>
</div>

{{end}}
