import { HttpClientModule } from '@angular/common/http';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { AceEditorModule } from 'ng2-ace-editor';
import { ToastrModule } from 'ngx-toastr';
import { AppComponent } from './app.component';
import { HolesComponent } from './holes/holes.component';
import { HomeComponent } from './home/home.component';
import { PlayviewComponent } from './playview/playview.component';
import { ProfileComponent } from './profile/profile.component';
import { LeaderboardComponent } from './leaderboard/leaderboard.component';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent
  },
  {
    path: 'play/:id',
    component: PlayviewComponent,
  },
  {
    path: 'holes',
    component: HolesComponent
  },
  {
    path: 'profile/:username',
    component: ProfileComponent
  }
];

@NgModule({
  declarations: [
    AppComponent,
    HolesComponent,
    PlayviewComponent,
    HomeComponent,
    ProfileComponent,
    LeaderboardComponent
  ],
  imports: [
    BrowserModule,
    ClarityModule,
    BrowserAnimationsModule,
    RouterModule.forRoot(
      routes,
    ),
    AceEditorModule,
    ToastrModule.forRoot(),
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class AppModule { }
