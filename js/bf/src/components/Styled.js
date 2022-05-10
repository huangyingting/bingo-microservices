import styled from 'styled-components';
import { keyframes, css } from "styled-components"
import { Col, Row } from 'react-bootstrap';
import React from 'react'
import { Fade } from "react-awesome-reveal";

const Section = styled.section`
  padding: 2em;
  background: ${props => props.color};
  display: block;
  vertical-align: middle;
`;

const animatedGradientH1 = keyframes`
  0%{background-position:0% 4%;}
  50%{background-position:100% 97%;}
  100%{background-position:0% 4%;}
`;

const AnimatedGradientH1 = styled.h1`
  background: linear-gradient(91.36deg, #ECA658 0%, #F391A6 13.02%, #E188C3 25.52%, #A58DE3 37.5%, #56ABEC 49.48%, #737EB7 63.02%, #C8638C 72.92%, #DD5D57 84.38%, #DF6C51 97.92%);
  background-size: 200% 200%;
  background-clip: text;
  animation: ${animatedGradientH1} 10s infinite ease both;
  font-size: calc(2rem + 1.5vw);
  font-weight: 700;
  overflow-wrap: break-word;
  text-align: center;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-background-clip: text;
`;

const FlatButton = styled.button`
  position: absolute;  
  top: 50%;
  transform: translateY(-50%);  
  cursor: pointer;
  display: inline-block;
  height: 64px;
  line-height: 56px;
  padding: 0 10px;
  color: white;
  width: 256px;
  background-color: transparent;
  bottom: 0px;
  left: 0px;
  right: 0px;
  margin: auto;
  border: solid 2px;
  //border-image: linear-gradient(to left, #743ad5 0%, #d53a9d 100%);
  border-image-slice: 1;
  font-size: 32px;
  font-weight: 600;
  clip-path: polygon(0 0, 12px 0, 12px 1px, 24px 1px, 24px 0, 100% 0, 100% 100%, 0 100%);
  &:hover{
    border: 0;
    background-color: rgba(365,365,365,0.5);
    cursor: pointer;
    color: #fff;
    opacity: 0.75;
    transition: .3s;
  }
  &:focus{
    outline: none;
    }   
`;

const GradientText = styled.div`
  background: #6a11cb;
  background: -webkit-linear-gradient(to right, rgba(119, 1, 103, 0.6), rgba(8, 0, 122, 0.6));
  background: linear-gradient(to right, rgba(119, 1, 103, 0.6), rgba(8, 0, 122, 0.6));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  overflow-wrap: break-word;
`
const Showcase = styled.div`
  margin: 0;
  height: 100vh;
  font-weight: 100;
  background: radial-gradient(#1640a2,#120c60);
  -webkit-overflow-Y: hidden;
  -moz-overflow-Y: hidden;
  -o-overflow-Y: hidden;
  overflow-y: hidden;
  -webkit-animation: loginFadeIn 1 1s ease-out;
  -moz-animation: loginFadeIn 1 1s ease-out;
  -o-animation: loginFadeIn 1 1s ease-out;
  animation: loginFadeIn 1 1s ease-out;
`

const showcaseFadeOut = keyframes`
  0%{opacity: 0;}
  30%{opacity: 1;}
  80%{opacity: .9;}
  100%{opacity: 0;}
`;

const showcaseFinalFade = keyframes`
  0%{opacity: 0;}
  30%{opacity: 1;}
  80%{opacity: .9;}
  100%{opacity: 1;}
`;


const showcaseFloatUp = keyframes`
  0%{top: 100vh; opacity: 0;}
  25%{opacity: 1;}
  50%{top: 0vh; opacity: .8;}
  75%{opacity: 1;}
  100%{top: -100vh; opacity: 0;}
`;

const ShowcaseHeadline = styled.p`
  position: absolute;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-weight: 700;
  color: white;
  font-size: 4em;
  opacity: 0;
  -webkit-animation: ${props => css`${showcaseFadeOut} 1 ${props.duration} ease-in`}; 
  -moz-animation: ${props => css`${showcaseFadeOut} 1 ${props.duration} ease-in`};
  -o-animation: ${props => css`${showcaseFadeOut} 1 ${props.duration} ease-in`};
  animation: ${props => css`${showcaseFadeOut} 1 ${props.duration} ease-in`};
  -webkit-animation-delay: ${props => props.delay};
  -moz-animation-delay: ${props => props.delay};
  -o-animation-delay: ${props => props.delay};
  animation-delay: ${props => props.delay};
`

const ShowcaseHeadlineFinal = styled.p`
  position: absolute;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-weight: 700;
  color: white;
  font-size: 4em;
  opacity: 0;
  -webkit-animation: ${props => css`${showcaseFinalFade} 1 ${props.duration} ease-in`}; 
  -moz-animation: ${props => css`${showcaseFinalFade} 1 ${props.duration} ease-in`};
  -o-animation: ${props => css`${showcaseFinalFade} 1 ${props.duration} ease-in`};
  animation: ${props => css`${showcaseFinalFade} 1 ${props.duration} ease-in`};
  -webkit-animation-fill-mode: forwards;
  -moz-animation-fill-mode: forwards;
  -o-animation-fill-mode: forwards;
  animation-fill-mode: forwards;  
  -webkit-animation-delay: ${props => props.delay};
  -moz-animation-delay: ${props => props.delay};
  -o-animation-delay: ${props => props.delay};
  animation-delay: ${props => props.delay};
`

const ShowcaseFloat = styled.div`
  position: absolute;
  width: 0px;
  opacity: .75;
  background-color: white;
  box-shadow: #e9f1f1 0px 0px 20px 2px;
  opacity: 0;
  top: 100vh;
  bottom: 0px;
  left: 0px;
  right: 0px;
  margin: auto;
  -webkit-animation: ${props => css`${showcaseFloatUp} ${props.duration} infinite linear`}; 
  -moz-animation: ${props => css`${showcaseFloatUp} ${props.duration} infinite linear`}; 
  -o-animation: ${props => css`${showcaseFloatUp} ${props.duration} infinite linear`}; 
  animation: ${props => css`${showcaseFloatUp} ${props.duration} infinite linear`}; 
  -webkit-transform: ${props => css`scale(${props.scale})`};
  -moz-transform: ${props => css`scale(${props.scale})`};
  -o-transform: ${props => css`scale(${props.scale})`};
  transform: ${props => css`scale(${props.scale})`};
  left: ${props => props.left};
`


const HeroHeader = (props) => {
  return (
    <Section color={props.color}>
      <Row>
        <Col>
          <AnimatedGradientH1 className='mb-4'>{props.title}</AnimatedGradientH1>
        </Col>
      </Row>
      <Row className='justify-content-center'>
        <Col className="col-12 col-md-6 col-lg-6 text-center align-self-center">
          <Fade direction="down">
            <h1 className="fw-bold mb-4">{props.subTitle}</h1>
          </Fade>
          <h3><span className='text-primary fw-bold'>{props.content.split(' ').shift() + " "}</span><span className="fw-light">{props.content.substr(props.content.indexOf(" ") + 1)}</span></h3>
        </Col>
        <Col className="col-12 col-md-5 col-lg-4 align-self-center">
          <img src={props.image} alt="hero header" />
        </Col>
      </Row>
    </Section>
  )
}

const Paragraph = (props) => {
  return (
    <Section color={props.color}>
      <div className={props.reverse? "d-flex justify-content-center flex-row-reverse" : "d-flex justify-content-center"}>
        <img src={props.image} className="w-25" alt="paragraph" />
        <div className="align-self-center w-50">
          <Fade direction="down">
            <h1 className="fw-bold mb-4">{props.title}</h1>
          </Fade>
          <h4 className="fw-light">{props.content}</h4>
        </div>
      </div>
    </Section>
  )
}

export {
  Section, AnimatedGradientH1, GradientText, Paragraph, HeroHeader,
  ShowcaseHeadline, ShowcaseHeadlineFinal, ShowcaseFloat,
  FlatButton, Showcase
};