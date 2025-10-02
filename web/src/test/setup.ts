import '@testing-library/jest-dom'
import React from 'react'

// Mock Clerk hooks for testing
vi.mock('@clerk/clerk-react', () => ({
  useUser: vi.fn(),
  useAuth: vi.fn(),
  SignIn: vi.fn(({ children }) => React.createElement('div', { 'data-testid': 'clerk-sign-in' }, children)),
  SignUp: vi.fn(({ children }) => React.createElement('div', { 'data-testid': 'clerk-sign-up' }, children)),
  SignedIn: vi.fn(),
  SignedOut: vi.fn(),
}))

// Mock React Router
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: vi.fn(),
    Navigate: vi.fn(),
    Link: vi.fn(({ children, to, ...props }) => 
      React.createElement('a', { href: to, ...props, 'data-testid': 'router-link' }, children)
    ),
  }
})

// Mock theme provider
vi.mock('@/components/theme-provider', () => ({
  useTheme: vi.fn(() => ({
    theme: 'light',
    setTheme: vi.fn(),
  })),
}))

// Mock mode toggle component
vi.mock('@/components/mode-toggle', () => ({
  ModeToggle: vi.fn(() => React.createElement('button', { 'data-testid': 'mode-toggle' }, 'Toggle Theme')),
}))