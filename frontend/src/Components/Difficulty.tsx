import { difficultyType } from "../Types";
import CSS from 'csstype';

export const Difficulty: React.FC<{difficulty: difficultyType; style?: CSS.Properties}> = ({difficulty, style}): JSX.Element => {
  let color = 'black';
  switch (difficulty) {
  case 'EASY':
    color = '#20C639'
    break;
  case 'IMPOSSIBLE':
    color = '#FF009D'
    break;
  case 'DIFFICULT':
  case 'HARD':
    color = '#E31717'
    break;
  case 'MEDIUM':
    color = '#EB9F45'
    break;
  }

  return (<p style={{color: color, ...style}}>{difficulty.toUpperCase()}</p>);
}
