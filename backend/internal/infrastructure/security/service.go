// Package security expõe a interface pública do JWTService.
// Mantemos o tipo concreto exportado para que o middleware possa recebê-lo sem
// criar dependência circular entre packages.
package security

// JWTServiceImpl é o nome do tipo concreto — exportado para injeção nos middlewares.
type JWTServiceImpl = jwtService
