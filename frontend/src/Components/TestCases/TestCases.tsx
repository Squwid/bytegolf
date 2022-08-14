import React from 'react';
import { SubmissionTests } from '../../Types';
import TestCase from './TestCase';

type Props = {
  holeID: string;
  tests: SubmissionTests;
  sizeStyle?: React.CSSProperties;
}

const TestCases: React.FC<Props> = (props) => {
  const keys = Object.keys(props.tests);

  return (
    <div style={props.sizeStyle} >
      <p style={{fontSize: '1.3rem', letterSpacing: '-0.09rem'}}>TEST CASES ({Object.keys(props.tests).length}):</p>
      {keys.map((t, i) => <TestCase holeID={props.holeID} test={props.tests[t]} testID={t} />)}
    </div>
  )
}

export default TestCases;