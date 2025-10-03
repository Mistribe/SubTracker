import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { useUser, SignIn } from '@clerk/clerk-react'
import { useNavigate, Link } from 'react-router-dom'
import { useTheme } from '@/components/theme-provider'
import SignInPage from '../SignInPage'

// Mock the dependencies
vi.mock('@clerk/clerk-react', () => ({
  useUser: vi.fn(),
  SignIn: vi.fn(() => <div data-testid="clerk-sign-in">Mocked SignIn</div>),
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

describe('SignInPage', () => {
  const mockNavigate = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    mockUseNavigate.mockReturnValue(mockNavigate)
    mockUseTheme.mockReturnValue({
      theme: 'light',
      setTheme: vi.fn(),
    })
  })

  it('renders sign-in page for unauthenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignInPage />)

    // Check AuthLayout content
    expect(screen.getByText('Welcome back')).toBeInTheDocument()
    expect(screen.getByText('Sign in to your account to continue')).toBeInTheDocument()

    // Check Clerk SignIn component is rendered
    expect(screen.getByTestId('clerk-sign-in')).toBeInTheDocument()

    // Check navigation link to sign-up
    expect(screen.getByText("Don't have an account?")).toBeInTheDocument()
    expect(screen.getByText('Sign up')).toBeInTheDocument()
    
    const signUpLink = screen.getByTestId('router-link')
    expect(signUpLink).toHaveAttribute('href', '/sign-up')
  })

  it('redirects authenticated users to dashboard', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: true,
      user: { id: 'user_123' } as any,
      isLoaded: true,
    })

    render(<SignInPage />)

    expect(mockNavigate).toHaveBeenCalledWith('/dashboard', { replace: true })
  })

  it('returns null for authenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: true,
      user: { id: 'user_123' } as any,
      isLoaded: true,
    })

    const { container } = render(<SignInPage />)

    expect(container.firstChild).toBeNull()
  })

  it('does not redirect unauthenticated users', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignInPage />)

    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('includes proper navigation link styling', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignInPage />)

    const signUpLink = screen.getByTestId('router-link')
    expect(signUpLink).toHaveClass(
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

    render(<SignInPage />)

    expect(screen.getByTestId('clerk-sign-in')).toBeInTheDocument()
    expect(mockUseTheme).toHaveBeenCalled()
  })

  it('renders with correct AuthLayout props', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignInPage />)

    // Verify AuthLayout receives correct props
    expect(screen.getByText('Welcome back')).toBeInTheDocument()
    expect(screen.getByText('Sign in to your account to continue')).toBeInTheDocument()
  })

  it('includes sign-up navigation text and link', () => {
    mockUseUser.mockReturnValue({
      isSignedIn: false,
      user: null,
      isLoaded: true,
    })

    render(<SignInPage />)

    // Check the complete navigation text
    expect(screen.getByText("Don't have an account?")).toBeInTheDocument()
    expect(screen.getByText('Sign up')).toBeInTheDocument()

    // Verify link points to correct route
    const signUpLink = screen.getByTestId('router-link')
    expect(signUpLink).toHaveAttribute('href', '/sign-up')
  })
})