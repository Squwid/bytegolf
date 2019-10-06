import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule, Routes } from '@angular/router';
import { AceEditorModule } from 'ng2-ace-editor';
import { AceConfigInterface, ACE_CONFIG } from 'ngx-ace-wrapper';
import { AppComponent } from './app.component';
import { HolesComponent } from './holes/holes.component';
import { NavbarComponent } from './navbar/navbar.component';
import { PlayViewComponent } from './play-view/play-view.component';

const appRoutes: Routes = [
  { path: 'holes', component: HolesComponent},
  { path: 'hole', component: PlayViewComponent}
];

const DEFAULT_ACE_CONFIG: AceConfigInterface = {

};


@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    HolesComponent,
    PlayViewComponent
  ],
  imports: [
    RouterModule.forRoot(appRoutes),
    BrowserModule,
    AceEditorModule
  ],
  providers: [
    {
      provide: ACE_CONFIG,
      useValue: DEFAULT_ACE_CONFIG
    }
],
  bootstrap: [AppComponent]
})
export class AppModule { }
