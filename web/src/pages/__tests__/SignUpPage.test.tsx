import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { useUser, SignUp } from '@clerk/clerk-react'
import { useNavigate, Link } from 'react-router-dom'
import { useTheme } from '@/components/theme-provider'
import SignUpPage from '../SignUpPage'

// Mock the dependencies
vi.mock('@clerk/clerk-react', () => ({
  useUser: vi.fn(),
  SignUp: vi.fn(() => <div data-testid="clerk-sign-up">Mocked SignUp</div>),
}))

vi.mock('react-router-dom', () => ({
  useNavigate: vi.fn(),
  Link: vi.fn(({ children, to, ...props }) => (
    <a href={to} {...props} data-testid="router-link">
      {children}
    </a>
  )),
}))

vi.mock('@/components/theme-provider', () => ({
  useTheme: vi.fn(),
}))

vi.mock('@/layouts/AuthLayout', () => ({
  AuthLayout: vi.fn(({ children, title, description }) => (
    <div data-testid="auth-layout">
      <h2>{title}</h2>
      <p>{description}</p>
      {children}
    </div>
  )),
}))

const mockUseUser = vi.mocked(useUser)
const mockUseNavigate = vi.mocked(useNavigate)
const mockUseTheme = vi.mocked(useTheme)

describe('SignUpPage', () => {
  const mockNavigate = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    mockUseNavigate.mockReturnValue(mockNavigate)
    mockUseTheme.mockReturnValue({
      theme: 'light',
      setTheme: vi.fn(),
    })
  })

  it('renders sign-up page for unauthenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    // Check AuthLayout content
    expect(screen.getByText('Create your account')).toBeInTheDocument()
    expect(screen.getByText('Sign up to start tracking your subscriptions')).toBeInTheDocument()

    // Check Clerk SignUp component is rendered
    expect(screen.getByTestId('clerk-sign-up')).toBeInTheDocument()

    // Check navigation link to sign-in
    expect(screen.getByText('Already have an account?')).toBeInTheDocument()
    expect(screen.getByText('Sign in')).toBeInTheDocument()
    
    const signInLink = screen.getByTestId('router-link')
    expect(signInLink).toHaveAttribute('href', '/sign-in')
  })

  it('redirects authenticated users to dashboard', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: true,
      user: { id: 'user_123' },
      isLoaded: true,
    })

    render(<SignUpPage />)

    expect(mockNavigate).toHaveBeenCalledWith('/dashboard', { replace: true })
  })

  it('returns null for authenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: true,
      user: { id: 'user_123' },
      isLoaded: true,
    })

    const { container } = render(<SignUpPage />)

    expect(container.firstChild).toBeNull()
  })

  it('does not redirect unauthenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('includes proper navigation link styling', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    const signInLink = screen.getByTestId('router-link')
    expect(signInLink).toHaveClass(
      'font-medium',
      'text-primary',
      'hover:text-primary/80',
      'transition-colors'
    )
  })

  it('handles theme changes correctly', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    // Test with dark theme
    mockUseTheme.mockReturnValue({
      theme: 'dark',
      setTheme: vi.fn(),
    })

    render(<SignUpPage />)

    expect(screen.getByTestId('clerk-sign-up')).toBeInTheDocument()
    expect(mockUseTheme).toHaveBeenCalled()
  })

  it('renders with correct AuthLayout props', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    // Verify AuthLayout receives correct props
    expect(screen.getByText('Create your account')).toBeInTheDocument()
    expect(screen.getByText('Sign up to start tracking your subscriptions')).toBeInTheDocument()
  })

  it('includes sign-in navigation text and link', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    // Check the complete navigation text
    expect(screen.getByText('Already have an account?')).toBeInTheDocument()
    expect(screen.getByText('Sign in')).toBeInTheDocument()

    // Verify link points to correct route
    const signInLink = screen.getByTestId('router-link')
    expect(signInLink).toHaveAttribute('href', '/sign-in')
  })

  it('configures Clerk appearance for theming', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignUpPage />)

    // Verify that the component renders without errors when theme is applied
    expect(screen.getByTestId('clerk-sign-up')).toBeInTheDocument()
  })

  it('handles different theme values', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    // Test with system theme
    mockUseTheme.mockReturnValue({
      theme: 'system',
      setTheme: vi.fn(),
    })

    render(<SignUpPage />)

    expect(screen.getByTestId('clerk-sign-up')).toBeInTheDocument()
    expect(mockUseTheme).toHaveBeenCalled()
  })
})