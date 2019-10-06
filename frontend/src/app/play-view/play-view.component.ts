import { Component, OnInit } from '@angular/core';
import 'brace';
import 'brace/mode/golang';
import 'brace/mode/javascript';
import 'brace/mode/python';
import 'brace/mode/text';
import 'brace/theme/cobalt';
import 'brace/theme/monokai';
import 'brace/theme/tomorrow';
import { Language } from '../shared/language.model';
import { Question } from '../shared/question.model';
import { Theme } from '../shared/theme.model';

// import { AceEditorModule } from 'ng2-ace-editor';
export const LANGUAGES: Language[] = [
  new Language('golang', 'https://upload.wikimedia.org/wikipedia/commons/2/23/Go_Logo_Aqua.svg', 'Golang 1.13', 'golang'),
  new Language('python', 'https://img.icons8.com/color/48/000000/python.png', 'Python 3.7', 'python'),
];

export const THEMES: Theme[] = [
  new Theme('cobalt'),
  new Theme('tomorrow'),
  new Theme('monokai')
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

  // current configuration settings
  mode = 'text';
  activeTheme = 'cobalt';
  disabled = false;
  content = 'Code here!';

  onSelectLanguage(language: Language) {
    console.log('selecing language: ', language.mode);
    this.mode = language.mode.toLowerCase();
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
      // dont submit if this is the case
      return;
    }
    console.log('mode: ' + this.mode, ' | code ' + this.content);
  }

  // variables regarding the editor

  constructor() { }

  ngOnInit() {
  }
}
