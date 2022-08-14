import React, { useEffect } from 'react';
import CSS from 'csstype';
import Nav from '../Components/Nav/Nav';
import { LightTextColor, PrimaryColor, readme, SecondaryColor, ThirdColor } from '../Globals';
import Logo from '../Logo/bytegolf_logo-half.png';
import Markdown from '../Components/Markdown';
import Button from '../Components/Button/Button';
import { useHistory } from 'react-router-dom';


const HomeScreen: React.FC = () => {
  const history = useHistory();

  const [ remoteMarkdown, setRemoteMarkdown ] = React.useState<string>('');

  useEffect(()=>{
    document.title = 'Bytegolf';

    fetch(readme)
      .then(async resp => setRemoteMarkdown((await resp.text())));
  },[]);

  return (
    <div>
      <Nav active={'home'}/>
      <div style={headerContainer}>
        <div style={headerContentContainer} >
          <img style={image} src={Logo} alt="Bytegolf Logo" />
          <p style={text}>BYTEGOLF</p>
          <p style={smallText}>SOLVE CODE PROBLEMS IN THE LEAST AMOUNT OF BYTES!</p>
          <div style={btnGroup}>
            <Button onPress={() => history.push(`/play`)} color={SecondaryColor} activeColor='white' fontSize='1.1rem' text="TEE OFF!"/>
            <Button onPress={() => window.open('https://github.com/Squwid/bytegolf', '_blank')} color={ThirdColor} activeColor='white' fontSize='1.1rem' text="SEE THE SOURCE CODE"/>
          </div>
        </div>
      </div>

      <div style={bodyContainer}>
        <Markdown markdown={remoteMarkdown} />
      </div>
    </div>
  )
}

const bodyContainer: CSS.Properties = {
  textAlign: 'left',
  width: '700px',
  margin: '0 auto',
  padding: '20px',
  height: 'auto',
  fontFamily: 'FiraCode'
}

const headerContainer: CSS.Properties = {
  width: '100%',
  minWidth: '980px',
  backgroundColor: PrimaryColor,
  minHeight: '26.25rem'
}

const headerContentContainer: CSS.Properties = {
  width: '40%',
  minWidth: '980px',
  margin: '0 auto',
  height: 'auto',

  
  display: 'flex',
  flexDirection: 'column',
  flexWrap: 'nowrap'
}

const image: CSS.Properties = {
  margin: '0 auto',
  maxHeight: '18.5rem',
  maxWidth: '18.5rem'
}

const text: CSS.Properties = {
  fontFamily: 'FiraCode',
  color: 'white',
  fontSize: '3.75rem',
  textAlign: 'center',
  margin: 0
}

const smallText: CSS.Properties = {
  fontFamily: 'FiraCode',
  color: LightTextColor,
  fontSize: '1.25rem',
  fontWeight: 'lighter',
  textAlign: 'center',
  margin: 0
}

const btnGroup: CSS.Properties = {
  width: 'auto',
  height: 'auto',
  margin: '0 auto',
  marginTop: '2rem',
  marginBottom: '0.625rem',

  display: 'flex',
  
}

export default HomeScreen;