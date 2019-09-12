import { Component } from '@angular/core';
import { LbUser } from '../shared/lbuser.model';

@Component({
    selector: 'app-lb',
    templateUrl: './lb.component.html',
    styleUrls: ['./lb.component.css']
})
export class LbComponent {
    leaders: LbUser[] = [
        new LbUser('bwhitelaw24', 9, 'Go'),
        new LbUser('collin123', 1, 'Python 3'),
        new LbUser('collin123', 3, 'Python 3'),
    ];

    constructor() {
    }


}
