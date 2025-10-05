import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { AuthLayout } from '../AuthLayout'

// Mock the ModeToggle component
vi.mock('@/components/mode-toggle', () => ({
  ModeToggle: () => <button data-testid="mode-toggle">Toggle Theme</button>,
}))

describe('AuthLayout', () => {
  const defaultProps = {
    title: 'Test Title',
    description: 'Test Description',
    children: <div data-testid="test-children">Test Content</div>,
  }

  it('renders the layout with correct structure', () => {
    render(<AuthLayout {...defaultProps} />)

    // Check header elements
    expect(screen.getByText('Recurrent Payment Tracker')).toBeInTheDocument()
    expect(screen.getByTestId('mode-toggle')).toBeInTheDocument()

    // Check main content
    expect(screen.getByText('Test Title')).toBeInTheDocument()
    expect(screen.getByText('Test Description')).toBeInTheDocument()
    expect(screen.getByTestId('test-children')).toBeInTheDocument()

    // Check footer
    const currentYear = new Date().getFullYear()
    expect(screen.getByText(`© ${currentYear} Recurrent Payment Tracker. All rights reserved.`)).toBeInTheDocument()
  })

  it('renders children content correctly', () => {
    const customChildren = <div data-testid="custom-content">Custom Auth Form</div>
    
    render(
      <AuthLayout title="Custom Title" description="Custom Description">
        {customChildren}
      </AuthLayout>
    )

    expect(screen.getByTestId('custom-content')).toBeInTheDocument()
    expect(screen.getByText('Custom Auth Form')).toBeInTheDocument()
  })

  it('displays the correct title and description', () => {
    const customProps = {
      title: 'Sign In',
      description: 'Welcome back to your account',
      children: <div>Auth form</div>,
    }

    render(<AuthLayout {...customProps} />)

    expect(screen.getByText('Sign In')).toBeInTheDocument()
    expect(screen.getByText('Welcome back to your account')).toBeInTheDocument()
  })

  it('has proper CSS classes for styling', () => {
    const { container } = render(<AuthLayout {...defaultProps} />)

    // Check main container has gradient background
    const mainContainer = container.firstChild as HTMLElement
    expect(mainContainer).toHaveClass('min-h-screen', 'bg-gradient-to-br', 'from-background', 'to-muted/20')

    // Check header structure
    const header = screen.getByRole('banner')
    expect(header).toHaveClass('container', 'mx-auto', 'p-4', 'flex', 'justify-between', 'items-center')

    // Check main content structure
    const main = screen.getByRole('main')
    expect(main).toHaveClass('container', 'mx-auto', 'px-4', 'py-8')
  })

  it('renders the application title as h1', () => {
    render(<AuthLayout {...defaultProps} />)

    const appTitle = screen.getByRole('heading', { level: 1 })
    expect(appTitle).toHaveTextContent('Recurrent Payment Tracker')
    expect(appTitle).toHaveClass('text-2xl', 'font-bold')
  })

  it('renders the page title as h2', () => {
    render(<AuthLayout {...defaultProps} />)

    const pageTitle = screen.getByRole('heading', { level: 2 })
    expect(pageTitle).toHaveTextContent('Test Title')
    expect(pageTitle).toHaveClass('text-3xl', 'font-bold', 'tracking-tight')
  })

  it('includes the ModeToggle component in header', () => {
    render(<AuthLayout {...defaultProps} />)

    const modeToggle = screen.getByTestId('mode-toggle')
    expect(modeToggle).toBeInTheDocument()
    
    // Check it's in the header
    const header = screen.getByRole('banner')
    expect(header).toContainElement(modeToggle)
  })
})