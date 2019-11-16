import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import 'brace';
import 'brace/mode/golang';
import 'brace/mode/javascript';
import 'brace/mode/python';
import 'brace/mode/text';
import 'brace/theme/cobalt';
import 'brace/theme/dracula';
import 'brace/theme/gob';
import 'brace/theme/monokai';
import 'brace/theme/terminal';
import 'brace/theme/textmate';
import { ToastrService } from 'ngx-toastr';
import { Execute } from '../shared/execute.model';
import { Language } from '../shared/language.model';
import { Question } from '../shared/question.model';
import { Theme } from '../shared/theme.model';

export const LANGUAGES: Language[] = [
  new Language('golang', 'https://upload.wikimedia.org/wikipedia/commons/2/23/Go_Logo_Aqua.svg', 'Golang 1.13', 'golang'),
  new Language('python3', 'https://img.icons8.com/color/48/000000/python.png', 'Python 3.7', 'python'),
];

const url = 'https://bytegolf.io';

export const THEMES: Theme[] = [
  new Theme('cobalt'),
  new Theme('tomorrow'),
  new Theme('monokai'),
  new Theme('dracula'),
  new Theme('textmate'),
  new Theme('terminal'),
  new Theme('gob')
];

@Component({
  selector: 'app-play-view',
  templateUrl: './play-view.component.html',
  styleUrls: ['./play-view.component.css']
})
export class PlayViewComponent implements OnInit {
  // all languages
  languages = LANGUAGES;
  themes = THEMES;

  hole = new Question('easy', 'Count to Ten!', 'This is a question where you have to count to 10 in code');
  submitDisabled = false;

  // current configuration settings
  mode = 'text';
  lang = '';
  activeTheme = 'cobalt';
  disabled = false;
  content = 'Code here!';

  onSelectLanguage(language: Language) {
    console.log('selecing language: ', language.mode);
    this.mode = language.mode;
    this.lang = language.language;
  }

  onChangeTheme(theme: Theme) {
    console.log('changing theme: ', theme);
    this.activeTheme = theme.name;
  }

  onClear() {
    this.content = '';
  }

  onSubmit() {
    if (this.content === '' || this.content === 'Code here!') {
      this.toastr.error('Code snippet cannot be blank', 'Error', {
        tapToDismiss: true,
      });
      return;
    }
    if (this.lang === '') {
      this.toastr.error('Language cannot be blank', 'Error', {
        tapToDismiss: true,
      });
      return;
    }
    this.toastr.info('', 'Sending Request', {
      tapToDismiss: true,
    });
    this.submitDisabled = true;

    // this is a good request create it and send here
    const exe = new Execute(this.content, this.lang, 'the only hole');

    this.http.post(url + '/compile', exe, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      })
    }).subscribe(responseData => {
      console.log(responseData);
      // todo: change this to not info when i can return a type from a request
      this.toastr.success(responseData.toString(), 'Success!', {
        tapToDismiss: true,
      });
      this.submitDisabled = false;
    });
  }


  // variables regarding the editor

  constructor(private http: HttpClient, private toastr: ToastrService) { }

  ngOnInit() {
  }
}
