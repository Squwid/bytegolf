import React, { useEffect } from 'react';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core';
import { makeStyles, Theme, withStyles } from '@material-ui/core/styles';
import { createStyles } from '@material-ui/styles';
import Nav from '../Components/Nav/Nav';
import { PrimaryColor } from '../Globals';
import './screen.css'
import { Difficulty } from '../Components/Difficulty';
import { useHistory } from 'react-router-dom';
import { ListHoles } from '../Store/Holes';
import { useQuery } from 'react-query';
import { BasicHole } from '../Types';
import { GetBestHoleScore } from '../Store/Subs';
import LoadingScreen from './PsuedoScreens/LoadingScreens';
import ErrorScreen from './PsuedoScreens/ErrorScreen';


type RowProps = {
  onClick: () => void;
  hole: BasicHole;
};

const Row: React.FC<RowProps> = (props) => {
  let score = '';
  const bestScore = useQuery(['BestScore', props.hole.ID], () => GetBestHoleScore(props.hole.ID));
  if (bestScore.isLoading || bestScore.isError || !bestScore.data) {
    score = '';
  } else if (bestScore.data === -1) {
    score = '-';
  } else {
    score = `${bestScore.data}`;
  } 


  return (
    <TRow key={props.hole.ID} onClick={props.onClick}>
      <TCell padding={'none'} style={{paddingLeft: '10px', paddingRight: '10px'}} component="th" scope="row">{props.hole.Name.toUpperCase()}</TCell>
      <TCell padding={'none'} style={{paddingLeft: '10px', paddingRight: '10px'}} align='right'>{score}</TCell>
      <TCell padding={'none'} style={{paddingLeft: '10px', paddingRight: '10px'}} align='right'><Difficulty difficulty={props.hole.Difficulty} /></TCell>
    </TRow>
  )
}


const HolesScreen: React.FC = () => {
  const classes = useStyles();
  const history = useHistory();

  useEffect(() => {document.title = `Bytegolf - Holes`}, [])

  const holes = useQuery('Holes', () =>  ListHoles());
  if (holes.isLoading) return (<LoadingScreen active='play' />);
  if (holes.isError || !holes.data) return (<ErrorScreen active='play' text={`${holes.error}`} />)

  return (
    <div>
      <Nav active='play' />
      <p className='screenTitle'>HOLES</p>
      <p className='screenText'>A LIST OF ALL ACTIVE HOLES</p>
      <div className='screenContainer'>
        <TableContainer >
          <Table className={classes.table}>
            <TableHead>
              <TableRow>
                <TCell >HOLE</TCell>
                <TCell align='right'>LOWEST SCORE</TCell>
                <TCell align='right'>DIFFICULTY</TCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {holes.data.map(hole => <Row hole={hole} onClick={() => history.push(`/play/${hole.ID}`)} />)}
            </TableBody>
          </Table>
        </TableContainer>
      </div>
    </div>
  );
}

const TCell = withStyles((theme: Theme) => createStyles({
  head: {
    backgroundColor: PrimaryColor,
    color: theme.palette.common.white,
    fontFamily: 'FiraCode',
    fontWeight: 'bold',
  },
  body: {
    fontSize: '1rem',
    fontFamily: 'FiraCode',
    fontWeight: 'lighter',
    letterSpacing: '-.09rem'
  }
})) (TableCell);

const TRow = withStyles((theme: Theme) => createStyles({
  root: {
    cursor: 'pointer',
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.action.hover,
    },
    '&:hover': {
      // textDecoration: 'underline',
      backgroundColor: 'lightgreen',
    },
    '&:active': {
      backgroundColor: PrimaryColor
    }
  }
})) (TableRow);

const useStyles = makeStyles({
  table: {
    width: '980px',
    border: '1px solid #CCC',
    marginBottom: '5rem',
    marginTop: '1rem'
  },
})

export default HolesScreen;