import { Row, Col } from 'react-bootstrap';
import { Fade } from "react-awesome-reveal";
import { AnimatedGradientH1 } from './Styled';

const ShortUrlHeader = (props) => {
  return (
    <Row className="justify-content-center align-items-center">
      <Col className="col-12 col-md-6 col-lg-6 text-center">
        <AnimatedGradientH1>Short links, big results</AnimatedGradientH1>
        <Fade direction="down">
          <h2><span className='text-primary fw-bold'>A Platform</span> with all features in one place. Shorten, brand, manage and track your links.</h2>          
        </Fade>
      </Col>
      <Col className="col-12 col-md-4 col-lg-3 text-center">
        <img src="/images/shorturl.svg" alt="shorturl" />
      </Col>
    </Row>
  )
}

export default ShortUrlHeader;