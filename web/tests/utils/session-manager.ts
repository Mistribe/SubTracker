/**
 * Session Management Utilities
 * 
 * Provides utilities for managing Clerk session tokens and persistence
 * across test runs and browser contexts
 */

import { Page, BrowserContext } from '@playwright/test';
import { ClerkTestUserManager, ClerkTestUser } from './clerk-user-manager';
import { getTestEnvironment } from '../config/test-env';

const testEnv = getTestEnvironment();

export interface SessionInfo {
  userId: string;
  sessionToken: string;
  expiresAt: number;
  isValid: boolean;
}

/**
 * Session Manager for Clerk authentication in tests
 */
export class SessionManager {
  private static sessions: Map<string, SessionInfo> = new Map();

  /**
   * Store session information for a user
   */
  static async storeSession(userId: string, sessionToken: string, expiresAt?: number): Promise<void> {
    const sessionInfo: SessionInfo = {
      userId,
      sessionToken,
      expiresAt: expiresAt || Date.now() + (24 * 60 * 60 * 1000), // Default 24 hours
      isValid: true,
    };

    this.sessions.set(userId, sessionInfo);
    
    // Also update in ClerkTestUserManager
    await ClerkTestUserManager.updateUserSessionToken(userId, sessionToken);
  }

  /**
   * Get session information for a user
   */
  static getSession(userId: string): SessionInfo | null {
    const session = this.sessions.get(userId);
    
    if (!session) return null;
    
    // Check if session is expired
    if (session.expiresAt < Date.now()) {
      session.isValid = false;
      this.sessions.delete(userId);
      return null;
    }
    
    return session;
  }

  /**
   * Check if user has valid session
   */
  static hasValidSession(userId: string): boolean {
    const session = this.getSession(userId);
    return session?.isValid || false;
  }

  /**
   * Invalidate session for a user
   */
  static invalidateSession(userId: string): void {
    this.sessions.delete(userId);
  }

  /**
   * Clear all sessions
   */
  static clearAllSessions(): void {
    this.sessions.clear();
  }

  /**
   * Extract session token from page context
   */
  static async extractSessionFromPage(page: Page): Promise<string | null> {
    try {
      // Extract Clerk session token from the page
      const sessionToken = await page.evaluate(() => {
        // Try to get session from Clerk's window object
        const clerk = (window as any).__clerk;
        if (clerk && clerk.session) {
          return clerk.session.id || clerk.session.token;
        }
        
        // Try to get from localStorage
        const clerkSession = localStorage.getItem('__clerk_session');
        if (clerkSession) {
          try {
            const parsed = JSON.parse(clerkSession);
            return parsed.token || parsed.id;
          } catch {
            return clerkSession;
          }
        }
        
        // Try to get from cookies
        const cookies = document.cookie.split(';');
        for (const cookie of cookies) {
          const [name, value] = cookie.trim().split('=');
          if (name.includes('clerk') && name.includes('session')) {
            return value;
          }
        }
        
        return null;
      });

      return sessionToken;
    } catch (error) {
      console.warn('Failed to extract session token from page:', error);
      return null;
    }
  }

  /**
   * Restore session in browser context
   */
  static async restoreSessionInContext(
    context: BrowserContext, 
    userId: string
  ): Promise<boolean> {
    const session = this.getSession(userId);
    if (!session) return false;

    try {
      // Add session token to browser context storage
      await context.addInitScript((sessionToken) => {
        // Set session in localStorage
        localStorage.setItem('__clerk_session', JSON.stringify({
          token: sessionToken,
          id: sessionToken,
        }));
        
        // Set session cookie
        document.cookie = `__clerk_session=${sessionToken}; path=/; max-age=86400`;
      }, session.sessionToken);

      return true;
    } catch (error) {
      console.warn('Failed to restore session in context:', error);
      return false;
    }
  }

  /**
   * Persist session from authenticated page
   */
  static async persistSessionFromPage(page: Page, userId: string): Promise<boolean> {
    const sessionToken = await this.extractSessionFromPage(page);
    
    if (!sessionToken) {
      console.warn('No session token found on page');
      return false;
    }

    await this.storeSession(userId, sessionToken);
    return true;
  }

  /**
   * Setup session persistence for test user
   */
  static async setupSessionPersistence(
    page: Page, 
    testUser: ClerkTestUser
  ): Promise<void> {
    try {
      // Wait for authentication to complete - try multiple approaches
      const sessionReady = await Promise.race([
        // Wait for Clerk session
        page.waitForFunction(() => {
          const clerk = (window as any).__clerk;
          return clerk && clerk.session;
        }, { timeout: 10000 }).then(() => true).catch(() => false),
        
        // Wait for URL to indicate successful auth
        page.waitForURL('/dashboard', { timeout: 10000 }).then(() => true).catch(() => false),
        
        // Wait for any auth-related localStorage
        page.waitForFunction(() => {
          return localStorage.getItem('__clerk_session') || 
                 document.cookie.includes('__session') ||
                 document.cookie.includes('clerk');
        }, { timeout: 10000 }).then(() => true).catch(() => false)
      ]);

      if (sessionReady) {
        // Extract and store session
        await this.persistSessionFromPage(page, testUser.id);
      } else {
        console.warn('Session persistence setup timed out, but continuing...');
      }
    } catch (error) {
      console.warn('Failed to setup session persistence:', error);
      // Don't fail the test, just log the warning
    }
  }

  /**
   * Restore session for test user
   */
  static async restoreUserSession(
    context: BrowserContext,
    testUser: ClerkTestUser
  ): Promise<boolean> {
    return await this.restoreSessionInContext(context, testUser.id);
  }

  /**
   * Check if session is still valid by testing page access
   */
  static async validateSessionOnPage(page: Page): Promise<boolean> {
    try {
      // Navigate to a protected route
      await page.goto('/dashboard');
      
      // Check if we're redirected to login
      await page.waitForLoadState('networkidle', { timeout: 5000 });
      
      const currentUrl = page.url();
      return !currentUrl.includes('/sign-in');
    } catch {
      return false;
    }
  }

  /**
   * Refresh session if needed
   */
  static async refreshSessionIfNeeded(
    page: Page, 
    testUser: ClerkTestUser
  ): Promise<boolean> {
    const isValid = await this.validateSessionOnPage(page);
    
    if (!isValid) {
      // Session expired, need to re-authenticate
      this.invalidateSession(testUser.id);
      return false;
    }
    
    // Update session token if it changed
    await this.persistSessionFromPage(page, testUser.id);
    return true;
  }

  /**
   * Get session expiration time
   */
  static getSessionExpiration(userId: string): number | null {
    const session = this.getSession(userId);
    return session?.expiresAt || null;
  }

  /**
   * Extend session expiration
   */
  static extendSession(userId: string, additionalTime: number = 24 * 60 * 60 * 1000): boolean {
    const session = this.sessions.get(userId);
    if (!session) return false;
    
    session.expiresAt = Date.now() + additionalTime;
    return true;
  }

  /**
   * Get all active sessions
   */
  static getActiveSessions(): SessionInfo[] {
    return Array.from(this.sessions.values()).filter(session => session.isValid);
  }

  /**
   * Cleanup expired sessions
   */
  static cleanupExpiredSessions(): void {
    const now = Date.now();
    for (const [userId, session] of this.sessions.entries()) {
      if (session.expiresAt < now) {
        this.sessions.delete(userId);
      }
    }
  }
}