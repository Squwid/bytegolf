import { Icon, makeStyles } from '@material-ui/core';
import React from 'react';
import { useQuery } from 'react-query';
import { PrimaryColor, SecondaryColor } from '../../Globals';
import { GetTest } from '../../Store/Holes';
import { SubmissionTestCase } from '../../Types';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

type Props = {
  test: SubmissionTestCase;
  testID: string;
  holeID: string;
}

const TestCase: React.FC<Props> = ({holeID, testID, test}) => {
  const [open, setOpen] = React.useState<boolean>(false);
  
  const classes = makeStyles({
    testCaseHeader: {
      fontFamily: 'FiraCode',
      fontSize: '1rem',
      border: `3px ${test.Correct ? PrimaryColor : SecondaryColor} solid`,
      borderBottom: `${open ? `0.1px lightgray solid` : `3px ${test.Correct ? PrimaryColor : SecondaryColor} solid`}`,
      padding: '0px',
      paddingLeft: '3px',
      paddingRight: '3px',
      marginTop: '5px',
      cursor: test.Hidden ? 'cursor' : 'pointer',
      letterSpacing: '-.09rem',
      backgroundColor: `${test.Correct ? '#BCFFC3' : '#FFCFC4'}`,
      '&:hover': {
        backgroundColor: test.Hidden ? `` : `${test.Correct ? 'lightgreen' : '#FFB6A5'}`
      },
      '&:active': {
        backgroundColor: test.Hidden ? `` : `${test.Correct ? '#2BDB7C' : '#FF9D87'}`
      }
    },
    testCaseBody: {
      border: `3px ${test.Correct ? PrimaryColor : SecondaryColor} solid`,
      borderTop: '0.1px lightgray solid',
      padding: '0px',
      paddingLeft: '3px',
      paddingRight: '3px',
      margin: '0px',
      // backgroundColor: 'lightcoral',
      
      display: 'flex',
      flexDirection: 'column',
    },
    testOutputMarkdown: {
      width: '90%',
      margin: '0 auto',
      padding: '0px',
      height: 'auto',
    },
  })();

  const onClick = () => {
    if (test.Hidden) return;
    setOpen(!open);
  }

  const testCase = useQuery(['TestCase', testID], () => GetTest(holeID, testID));
  if (testCase.isError){
    console.error(testCase.error);
    return (<></>);
  }
  if (!testCase.data && !testCase.isLoading) return (<></>);

  // Logic to choose what to display
  let name = '...';
  if (!testCase.isLoading && !testCase.data) name = '';
  else if (!testCase.isLoading && testCase.data) {
    name = testCase.data.Name;
  }

  return (
    <>
    <div className={classes.testCaseHeader} onClick={onClick}>
      <div style={{display: 'flex', flexDirection: 'row', flexWrap: 'nowrap', justifyContent: 'space-between', alignItems: 'center'}}>
        <p>{name}</p>
        {!testCase.data?.Hidden && (
          <div>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </div>
        )}

        {testCase.data?.Hidden && (
          <p>(HIDDEN)</p>
        )}
      </div>
    </div>

    {open && (
      <div className={classes.testCaseBody}>
        <p style={{padding: 10, margin: 0}}>ID: {testID}</p>

        {testCase.data?.Input && <><p style={{padding: 10, margin: 0}}>INPUT:</p>
        <div className={classes.testOutputMarkdown}>
          <p style={{borderRadius: '5px', backgroundColor: '#f0f0f0', padding: '10px', minHeight: '10px'}}>{testCase.data?.Input}</p>
        </div></>}

        <p style={{padding: 10, margin: 0}}>OUTPUT:</p>
        <div className={classes.testOutputMarkdown}>
          <p style={{borderRadius: '5px', backgroundColor: '#f0f0f0', padding: '10px', minHeight: '10px'}}>{test.Output}</p>
        </div>
      </div>
    )}
    </>

  );
}

export default TestCase;