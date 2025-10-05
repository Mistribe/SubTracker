import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { SignedIn, SignedOut } from '@clerk/clerk-react'
import { Navigate } from 'react-router-dom'
import { ProtectedRoute } from '../ProtectedRoute'

// Mock Clerk components
vi.mock('@clerk/clerk-react', () => ({
  SignedIn: vi.fn(),
  SignedOut: vi.fn(),
}))

// Mock React Router Navigate
vi.mock('react-router-dom', () => ({
  Navigate: vi.fn(),
}))

const mockSignedIn = vi.mocked(SignedIn)
const mockSignedOut = vi.mocked(SignedOut)
const mockNavigate = vi.mocked(Navigate)

describe('ProtectedRoute', () => {
  const TestChildren = () => <div data-testid="protected-content">Protected Content</div>

  beforeEach(() => {
    vi.clearAllMocks()
    
    // Default mock implementations
    mockSignedIn.mockImplementation(({ children }) => <>{children}</>)
    mockSignedOut.mockImplementation(({ children }) => <>{children}</>)
    mockNavigate.mockImplementation(() => null)
  })

  it('renders children when user is signed in', () => {
    // Mock SignedIn to render children, SignedOut to render nothing
    mockSignedIn.mockImplementation(({ children }) => <>{children}</>)
    mockSignedOut.mockImplementation(() => null)

    render(
      <ProtectedRoute>
        <TestChildren />
      </ProtectedRoute>
    )

    expect(screen.getByTestId('protected-content')).toBeInTheDocument()
    expect(screen.queryByTestId('navigate-component')).not.toBeInTheDocument()
  })

  it('redirects to sign-in when user is signed out', () => {
    // Mock SignedIn to render nothing, SignedOut to render children
    mockSignedIn.mockImplementation(() => null)
    mockSignedOut.mockImplementation(({ children }) => <>{children}</>)

    render(
      <ProtectedRoute>
        <TestChildren />
      </ProtectedRoute>
    )

    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument()
    
    // Verify Navigate is called with correct props
    expect(mockNavigate).toHaveBeenCalledWith({ to: '/sign-in', replace: true }, undefined)
  })

  it('uses SignedIn and SignedOut components correctly', () => {
    render(
      <ProtectedRoute>
        <TestChildren />
      </ProtectedRoute>
    )

    // Verify that both Clerk components are used
    expect(mockSignedIn).toHaveBeenCalledWith({ children: <TestChildren /> }, undefined)
    expect(mockSignedOut).toHaveBeenCalled()
  })

  it('passes replace prop to Navigate component', () => {
    // Mock SignedOut to render the Navigate component
    mockSignedIn.mockImplementation(() => null)
    mockSignedOut.mockImplementation(({ children }) => <>{children}</>)

    render(
      <ProtectedRoute>
        <TestChildren />
      </ProtectedRoute>
    )

    expect(mockNavigate).toHaveBeenCalledWith(
      expect.objectContaining({
        to: '/sign-in',
        replace: true,
      }),
      undefined
    )
  })

  it('handles multiple children correctly', () => {
    mockSignedIn.mockImplementation(({ children }) => <>{children}</>)
    mockSignedOut.mockImplementation(() => null)

    render(
      <ProtectedRoute>
        <div data-testid="child-1">Child 1</div>
        <div data-testid="child-2">Child 2</div>
      </ProtectedRoute>
    )

    expect(screen.getByTestId('child-1')).toBeInTheDocument()
    expect(screen.getByTestId('child-2')).toBeInTheDocument()
  })

  it('redirects to correct sign-in route', () => {
    mockSignedIn.mockImplementation(() => null)
    mockSignedOut.mockImplementation(({ children }) => <>{children}</>)

    render(
      <ProtectedRoute>
        <TestChildren />
      </ProtectedRoute>
    )

    // Verify the redirect goes to the correct route
    expect(mockNavigate).toHaveBeenCalledWith(
      expect.objectContaining({
        to: '/sign-in',
      }),
      undefined
    )
  })
})