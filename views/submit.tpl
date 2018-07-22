{{define "submit"}}

<!-- Submission -->
<div class="content">
    <div class="contentspace"></div>
    <div class="golfcontainer">
        <h2>Submission</h2>
        <div class="golfoptions">
            <form method="POST" enctype="multipart/form-data">
                <que>Current Game &nbsp;</que>
                <br> {{.GameName}}
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
                <input type="submit" value="Tee Off" style="height:25px; width:80px">

            </form>
        </div>
    </div>

</div>

{{end}}
