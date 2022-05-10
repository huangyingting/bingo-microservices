import { Card, Row, Col, Container } from "react-bootstrap";
import { Flip } from "react-awesome-reveal";
import { Section, HeroHeader } from './Styled';

function About() {
  return (
    <>
      <HeroHeader color="#f7f5eb" title="Bingo"
        subTitle="To empower every person and every organization on the planet to achieve more"
        content="Our culture is centered on embracing a growth mindset, a theme of inspiring excellence, and encouraging teams and leaders to bring their best each day. In doing so, we create life-changing innovations that impact billions of lives around the world."
        image="/images/team.svg"
      />
      <Section color="#fff">
        <Container>
          <Row className="row-cols-1 row-cols-md-2 row-cols-xl-4 gy-4">
            <Col>
              <Flip direction="horizontal">
                <Card className='shadow h-100'>
                  <Card.Body>
                    <Card.Title className="text-center fw-bold">Our company</Card.Title>
                    <Card.Text>Founded in 1975.Stay informed about us - from company facts and news to our worldwide locations and more.</Card.Text>
                  </Card.Body>
                </Card>
              </Flip>
            </Col>
            <Col>
              <Flip direction="horizontal">
                <Card className='shadow h-100'>
                  <Card.Body>
                    <Card.Title className="text-center fw-bold">Who are we</Card.Title>
                    <Card.Text>Get to know some of our people, explore engaging stories, and meet the leaders who shape our vision.</Card.Text>
                  </Card.Body>
                </Card>
              </Flip>
            </Col>
            <Col>
              <Flip direction="horizontal">
                <Card className='shadow h-100'>
                  <Card.Body>
                    <Card.Title className="text-center fw-bold">What we value</Card.Title>
                    <Card.Text>See how we utilize technology to build platforms and resources to help make a lasting positive impact.</Card.Text>
                  </Card.Body>
                </Card>
              </Flip>
            </Col>
            <Col>
              <Flip direction="horizontal">
                <Card className='shadow h-100'>
                  <Card.Body>
                    <Card.Title className="text-center fw-bold">Contact us</Card.Title>
                    <Card.Text>Meet our brilliant and knowledgeable cloud solution architect team. Get in touch. We’re here to help.</Card.Text>
                  </Card.Body>
                </Card>
              </Flip>
            </Col>
          </Row>
        </Container>
      </Section>
    </>
  )
}
export default About;