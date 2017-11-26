﻿using System;
using System.Net.Http;
using System.Threading.Tasks;
using Restaurant.Abstractions.Api;
using Restaurant.Abstractions.Providers;
using Restaurant.Common.DataTransferObjects;

namespace Restaurant.Core.Providers
{
    public class AuthenticationProvider : IAuthenticationProvider
    {
        private readonly ITokenProvider _tokenProvider;
        private readonly IAccountApi _accountApi;

        public AuthenticationProvider(
            ITokenProvider tokenProvider,
            IAccountApi accountApi)
        {
            _tokenProvider = tokenProvider;
            _accountApi = accountApi;
        }

        public Task<TokenResponse> Login(LoginDto loginDto)
        {
            return _tokenProvider.RequestResourceOwnerPasswordAsync(loginDto.Login, loginDto.Password);
        }

        public Task<HttpResponseMessage> Register(RegisterDto registerDto)
        {
            return _accountApi.Register(registerDto);
        }

	    public Task<TokenResponse> RefreshToken(string refreshToken)
	    {
		    throw new NotImplementedException();
	    }

	    public Task<object> LogOut()
        {
            return _accountApi.LogOut();
        }
    }
}
