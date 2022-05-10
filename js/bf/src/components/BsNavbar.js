import './BsNavbar.css'
import { useState, useEffect } from 'react';
import { useLocation, Link } from 'react-router-dom';
import { Navbar, Container, Nav, Button } from 'react-bootstrap'
import { useMsal, useIsAuthenticated } from "@azure/msal-react";
import { loginRequest } from "../AAD";
import { AiOutlineLink } from 'react-icons/ai';

export default function BsNavbar() {
  let location = useLocation()
  const [path, setPath] = useState(null);
  const { instance, accounts } = useMsal();
  const isAuthenticated = useIsAuthenticated();
  const name = accounts[0] && accounts[0].name;

  useEffect(() => {
    setPath(location.pathname);
  }, [location]);

  function handleLogin() {
    instance.loginPopup(loginRequest).catch(e => {
      console.error(e);
    });
  }

  function handleLogout() {
    instance.logoutPopup().catch(e => {
      console.error(e);
    });
  }

  return (
    <Navbar collapseOnSelect expand="md" bg="white" fixed="top" sticky='top' className='shadow'>
      <Container fluid>
        <Navbar.Brand href="/">
          <h3 className='fw-bold'><AiOutlineLink />Bingo</h3>           
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="collapse" />
        <Navbar.Collapse id="collapse">
          <Nav className="me-auto">
            <Nav.Link as={Link} to='/' href='/' active={path === '/'}>
            <h5 className='fw-bold'>Home</h5>
            </Nav.Link>
            <Nav.Link as={Link} to='/pages/stats' href='/pages/stats' active={path === '/pages/stats'}>
            <h5 className='fw-bold'>System</h5>
            </Nav.Link>
            <Nav.Link as={Link} to='/pages/shorturl' href='/pages/shorturl' active={path === '/pages/shorturl'}>
            <h5 className='fw-bold'>Shorten Links</h5>
            </Nav.Link>
            <Nav.Link as={Link} to='/pages/about' href='/pages/about' active={path === '/pages/about'}>
              <h5 className='fw-bold'>About</h5>
            </Nav.Link>
          </Nav>
          {isAuthenticated && <div className='px-4 large text-truncate'><strong>Hello</strong> {name}</div>}
          {isAuthenticated ?
            <Button variant='outline-primary' className='rounded px-4' onClick={() => handleLogout()}>Logout</Button> :
            <Button variant='outline-primary' className='rounded px-4' onClick={() => handleLogin()}>Login</Button>
          }
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}
