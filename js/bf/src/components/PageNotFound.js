import { Section } from './Styled';
import { AnimatedGradientH1 } from './Styled';
import { useNavigate } from 'react-router-dom'
import { Button } from 'react-bootstrap';


function PageNotFound() {
  const navigate = useNavigate()

  return (
    <Section color="#f7f5eb">
      <AnimatedGradientH1>404 - Page not found!</AnimatedGradientH1>
      <div className="d-flex justify-content-center m-4">
        <h3><b className='text-primary'>Oops!</b> Itseems like you're lost, let me help you find your way back home! :)<br />
          The following problems could be the case:<br />
          • The link is broken<br />
          • This page never existed<br />
          • This page has been removed</h3>
      </div>
      <div className="text-center">
        <Button size="lg" variant='outline-primary' className='rounded px-4' onClick={() => { navigate('/') }}>
          Return
        </Button>
      </div>
    </Section>
  )
}
export default PageNotFound;