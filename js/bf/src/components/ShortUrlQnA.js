import { Accordion } from 'react-bootstrap';
import { Section } from './Styled';
import { Col, Row } from 'react-bootstrap';

const ShortUrlQnA = (props) => {
  return (
    <Section color="#fff">
      <Row className="justify-content-center align-items-center">
        <Col className="col-12 col-md-10 col-lg-9">
          <h1 className='fw-bold text-center'>Frequently asked questions</h1>
          <br />
          <br />
          <Accordion>
            <Accordion.Item eventKey="0">
              <Accordion.Header><h5 className='fw-bold'>What is a URL Shortener?</h5></Accordion.Header>
              <Accordion.Body>
                A URL shortener, also known as a link shortener, seems like a simple tool, but it is a service that can have a dramatic impact on your marketing efforts.
                Link shorteners work by transforming any long URL into a shorter, more readable link. When a user clicks the shortened version, they're automatically forwarded to the destination URL.
                Think of a short URL as a more descriptive and memorable nickname for your long webpage address. You can, for example, use a short URL so people will have a good idea about where your link will lead before they click it.
                If you're contributing content to the online world, you need a URL shortener.
                Make your URLs stand out with our easy to use free link shortener above.
              </Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="1">
              <Accordion.Header><h5 className='fw-bold'>Benefits of a Short URL</h5></Accordion.Header>
              <Accordion.Body>
                How many people can even remember a long web address, especially if it has tons of characters and symbols? A short URL can make your link more memorable. Not only does it allow people to easily recall and share your link with others, it can also dramatically improve traffic to your content.
                On a more practical side, a short URL is also easier to incorporate into your collateral - whether you're looking to engage with your customers offline or online.
                Bingo is the best URL shortener for everyone, from influencers to small brands to large enterprises, who are looking for a simple way to create, track and manage their links.
              </Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="2">
              <Accordion.Header><h5 className='fw-bold'>Manage, Monitor, and Measure</h5></Accordion.Header>
              <Accordion.Body>
                Link Management is more than shortening a link and branding is just the beginning.
                Bingo creates an environment to manage every touch, monitor for accuracy, and provides analytical insights to measure value and performance.  If you share a link, protect it with Bingo.
              </Accordion.Body>
            </Accordion.Item>
          </Accordion>
        </Col>
      </Row>
    </Section>
  )
}

export default ShortUrlQnA;