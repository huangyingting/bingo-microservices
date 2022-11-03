import { useMsal } from "@azure/msal-react";
import { loginRequest } from "../AAD";
import { FlatButton, Showcase, ShowcaseHeadline, ShowcaseHeadlineFinal, ShowcaseFloat } from "./Styled"

function Unauthenticated() {
  const { instance } = useMsal();
  function handleLogin() {
    instance.loginPopup(loginRequest).catch(e => {
      console.error(e);
    });
  }

  return (
    <Showcase>
      <ShowcaseHeadline duration="5s" delay="0s">A Cloud-native application</ShowcaseHeadline>
      <ShowcaseHeadline duration="5s" delay="6s">Designed</ShowcaseHeadline>
      <ShowcaseHeadline duration="5s" delay="12s">For resilience</ShowcaseHeadline>
      <ShowcaseHeadline duration="5s" delay="17s">Scalable and reliable all the time</ShowcaseHeadline>
      <ShowcaseHeadlineFinal duration="5s" delay="22s">Welcome to Bingo</ShowcaseHeadlineFinal>
      <FlatButton onClick={() => handleLogin()}>Login</FlatButton>
      <ShowcaseFloat duration="4s" scale="1.0" left="0%"></ShowcaseFloat>
      <ShowcaseFloat duration="7s" scale="1.6" left="15%"></ShowcaseFloat>
      <ShowcaseFloat duration="2.5s" scale=".5" left="-15%"></ShowcaseFloat>
      <ShowcaseFloat duration="4.5s" scale="1.2" left="-34%"></ShowcaseFloat>
      <ShowcaseFloat duration="8s" scale="2.2" left="-57%"></ShowcaseFloat>
      <ShowcaseFloat duration="3s" scale="0.8" left="-81%"></ShowcaseFloat>
      <ShowcaseFloat duration="5.3s" scale="3.2" left="37%"></ShowcaseFloat>
      <ShowcaseFloat duration="4.7s" scale="1.7" left="62%"></ShowcaseFloat>
      <ShowcaseFloat duration="4.1s" scale=".9" left="85%"></ShowcaseFloat>
    </Showcase>
  )
}
export default Unauthenticated;